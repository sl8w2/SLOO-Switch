package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	iniName     = "GameUserSettings.ini"
	keepBackups = 5
	version     = "1.0"
)

var keyRe = regexp.MustCompile(`^[a-f0-9]{16}$`)

type apiError struct {
	Code string
	Arg  string
}

func (e *apiError) Error() string { return e.Code }

func fail(code string, arg ...string) error {
	e := &apiError{Code: code}
	if len(arg) > 0 {
		e.Arg = arg[0]
	}
	return e
}

type Account struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	ID    string `json:"id"`
	Saved string `json:"saved"`
}

type AccountView struct {
	Account
	Active bool `json:"active"`
}

type Current struct {
	ID       string `json:"id"`
	LoggedIn bool   `json:"loggedIn"`
	Matched  string `json:"matched"`
}

type View struct {
	Accounts    []AccountView `json:"accounts"`
	Current     Current       `json:"current"`
	EpicRunning bool          `json:"epicRunning"`
	ExeFound    bool          `json:"exeFound"`
	ExeAuto     string        `json:"exeAuto"`
	ExePath     string        `json:"exePath"`
	DataPath    string        `json:"dataPath"`
	AutoSave    bool          `json:"autoSave"`
	Version     string        `json:"version"`
}

type SwitchResult struct {
	Launched bool   `json:"launched"`
	Name     string `json:"name"`
}

type settings struct {
	ExePath     string `json:"exePath"`
	AutoSaveOff bool   `json:"autoSaveOff"`
	Lang        string `json:"lang"`
}

type state struct {
	Active string `json:"active"`
	Ignore string `json:"ignore"`
}

type Store struct {
	mu     sync.Mutex
	dir    string
	p      platform
	ini    string
	caches []string
}

func newStore(dir string, p platform) *Store {
	local := os.Getenv("LOCALAPPDATA")
	return &Store{
		dir: dir,
		p:   p,
		ini: filepath.Join(local, "EpicGamesLauncher", "Saved", "Config", "WindowsEditor", iniName),
		caches: []string{
			filepath.Join(local, "Epic Games", "Epic Online Services", "UI Helper", "Cache", "Cache"),
			filepath.Join(local, "Epic Games", "Epic Online Services", "UI Helper", "Cache", "GPUCache"),
			filepath.Join(local, "Epic Games", "EOSOverlay", "BrowserCache", "Cache"),
		},
	}
}

func readJSON(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func writeFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return writeFile(path, b)
}

func cleanName(n string) string {
	n = strings.TrimSpace(n)
	if n == "" || utf8.RuneCountInString(n) > 40 {
		return ""
	}
	for _, r := range n {
		if unicode.IsControl(r) {
			return ""
		}
	}
	return n
}

func newKey() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func hasLogin(data []byte) bool {
	return bytes.Contains(bytes.ToLower(data), []byte("rememberme"))
}

func (s *Store) accDir(key string) string { return filepath.Join(s.dir, "accounts", key) }

func (s *Store) accounts() []Account {
	entries, _ := os.ReadDir(filepath.Join(s.dir, "accounts"))
	out := []Account{}
	for _, e := range entries {
		if !e.IsDir() || !keyRe.MatchString(e.Name()) {
			continue
		}
		var a Account
		if readJSON(filepath.Join(s.accDir(e.Name()), "meta.json"), &a) == nil {
			a.Key = e.Name()
			out = append(out, a)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out
}

func (s *Store) loadState() state {
	var st state
	readJSON(filepath.Join(s.dir, "state.json"), &st)
	return st
}

func (s *Store) saveState(st state) {
	writeJSON(filepath.Join(s.dir, "state.json"), st)
}

func (s *Store) settings() settings {
	var st settings
	readJSON(filepath.Join(s.dir, "settings.json"), &st)
	return st
}

func (s *Store) exePath() string {
	if p := s.settings().ExePath; p != "" {
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			return p
		}
	}
	for _, env := range []string{"ProgramFiles(x86)", "ProgramFiles"} {
		base := os.Getenv(env)
		if base == "" {
			continue
		}
		for _, arch := range []string{"Win32", "Win64"} {
			p := filepath.Join(base, "Epic Games", "Launcher", "Portal", "Binaries", arch, "EpicGamesLauncher.exe")
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
				return p
			}
		}
	}
	return ""
}

func (s *Store) View() View {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.p.ReadAccountID()
	data, err := os.ReadFile(s.ini)
	loggedIn := err == nil && hasLogin(data)
	v := View{Accounts: []AccountView{}, Version: version, EpicRunning: s.p.EpicRunning()}
	matched := ""
	for _, a := range s.accounts() {
		act := loggedIn && id != "" && a.ID == id
		if act {
			matched = a.Key
		}
		v.Accounts = append(v.Accounts, AccountView{Account: a, Active: act})
	}
	v.Current = Current{ID: id, LoggedIn: loggedIn, Matched: matched}
	v.ExeAuto = s.exePath()
	v.ExeFound = v.ExeAuto != ""
	v.ExePath = s.settings().ExePath
	v.DataPath = s.dir
	v.AutoSave = !s.settings().AutoSaveOff
	return v
}

func (s *Store) Save(name string) (Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	name = cleanName(name)
	if name == "" {
		return Account{}, fail("bad_name")
	}
	data, err := os.ReadFile(s.ini)
	if err != nil || !hasLogin(data) {
		return Account{}, fail("no_login")
	}
	id := s.p.ReadAccountID()
	if id == "" {
		return Account{}, fail("no_id")
	}
	for _, a := range s.accounts() {
		if strings.EqualFold(a.Name, name) {
			return Account{}, fail("name_taken", a.Name)
		}
		if a.ID == id {
			return Account{}, fail("id_taken", a.Name)
		}
	}
	a := Account{Key: newKey(), Name: name, ID: id, Saved: time.Now().Format(time.RFC3339)}
	if err := writeFile(filepath.Join(s.accDir(a.Key), iniName), data); err != nil {
		return Account{}, fail("internal")
	}
	if err := writeJSON(filepath.Join(s.accDir(a.Key), "meta.json"), a); err != nil {
		return Account{}, fail("internal")
	}
	st := s.loadState()
	st.Active = a.Key
	s.saveState(st)
	return a, nil
}

func (s *Store) refreshLocked() string {
	data, err := os.ReadFile(s.ini)
	if err != nil || !hasLogin(data) {
		return ""
	}
	id := s.p.ReadAccountID()
	if id == "" {
		return ""
	}
	var match *Account
	accs := s.accounts()
	for i := range accs {
		if accs[i].ID == id {
			match = &accs[i]
			break
		}
	}
	if match == nil {
		return ""
	}
	snap := filepath.Join(s.accDir(match.Key), iniName)
	old, _ := os.ReadFile(snap)
	if !bytes.Equal(old, data) {
		if writeFile(snap, data) == nil {
			match.Saved = time.Now().Format(time.RFC3339)
			writeJSON(filepath.Join(s.accDir(match.Key), "meta.json"), match)
		}
	}
	if st := s.loadState(); st.Active != match.Key {
		st.Active = match.Key
		s.saveState(st)
	}
	return match.Key
}

func (s *Store) Refresh() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refreshLocked()
}

func (s *Store) backupLocked() {
	data, err := os.ReadFile(s.ini)
	if err != nil {
		return
	}
	root := filepath.Join(s.dir, "backups")
	dest := filepath.Join(root, time.Now().Format("20060102-150405.000000000"))
	if writeFile(filepath.Join(dest, iniName), data) != nil {
		return
	}
	writeFile(filepath.Join(dest, "id.txt"), []byte(s.p.ReadAccountID()))
	entries, _ := os.ReadDir(root)
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for len(names) > keepBackups {
		os.RemoveAll(filepath.Join(root, names[0]))
		names = names[1:]
	}
}

func (s *Store) find(key string) (Account, bool) {
	if !keyRe.MatchString(key) {
		return Account{}, false
	}
	for _, a := range s.accounts() {
		if a.Key == key {
			return a, true
		}
	}
	return Account{}, false
}

func (s *Store) Switch(key string) (SwitchResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	acc, ok := s.find(key)
	if !ok {
		return SwitchResult{}, fail("not_found")
	}
	data, err := os.ReadFile(filepath.Join(s.accDir(key), iniName))
	if err != nil {
		return SwitchResult{}, fail("snapshot_missing")
	}
	if err := s.p.KillEpic(); err != nil {
		return SwitchResult{}, fail("kill_failed")
	}
	s.refreshLocked()
	s.backupLocked()
	if err := writeFile(s.ini, data); err != nil {
		return SwitchResult{}, fail("internal")
	}
	if acc.ID != "" {
		s.p.WriteAccountID(acc.ID)
	}
	for _, c := range s.caches {
		os.RemoveAll(c)
	}
	st := s.loadState()
	st.Active = key
	s.saveState(st)
	res := SwitchResult{Name: acc.Name}
	if exe := s.exePath(); exe != "" && s.p.Launch(exe) == nil {
		res.Launched = true
	}
	return res, nil
}

func (s *Store) NewLogin() (SwitchResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if data, err := os.ReadFile(s.ini); err == nil && hasLogin(data) {
		id := s.p.ReadAccountID()
		known := false
		for _, a := range s.accounts() {
			if id != "" && a.ID == id {
				known = true
			}
		}
		if !known {
			return SwitchResult{}, fail("unsaved_current")
		}
	}
	if err := s.p.KillEpic(); err != nil {
		return SwitchResult{}, fail("kill_failed")
	}
	s.refreshLocked()
	s.backupLocked()
	os.Remove(s.ini)
	for _, c := range s.caches {
		os.RemoveAll(c)
	}
	st := s.loadState()
	st.Active = ""
	s.saveState(st)
	res := SwitchResult{}
	if exe := s.exePath(); exe != "" && s.p.Launch(exe) == nil {
		res.Launched = true
	}
	return res, nil
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	acc, ok := s.find(key)
	if !ok {
		return fail("not_found")
	}
	if err := os.RemoveAll(s.accDir(key)); err != nil {
		return fail("internal")
	}
	st := s.loadState()
	if st.Active == key {
		st.Active = ""
	}
	if data, err := os.ReadFile(s.ini); err == nil && hasLogin(data) && acc.ID != "" && s.p.ReadAccountID() == acc.ID {
		st.Ignore = acc.ID
	}
	s.saveState(st)
	return nil
}

func (s *Store) updateSettings(fn func(*settings)) error {
	st := s.settings()
	fn(&st)
	return writeJSON(filepath.Join(s.dir, "settings.json"), st)
}

func (s *Store) SetExe(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	path = strings.TrimSpace(strings.Trim(strings.TrimSpace(path), `"`))
	if path != "" {
		fi, err := os.Stat(path)
		if err != nil || fi.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".exe") {
			return fail("bad_path")
		}
	}
	if s.updateSettings(func(st *settings) { st.ExePath = path }) != nil {
		return fail("internal")
	}
	return nil
}

func (s *Store) SetAutoSave(on bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.updateSettings(func(st *settings) { st.AutoSaveOff = !on }) != nil {
		return fail("internal")
	}
	return nil
}

func (s *Store) SetLang(lang string) {
	if lang != "ar" && lang != "en" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.updateSettings(func(st *settings) { st.Lang = lang })
}

func (s *Store) Rename(key, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	name = cleanName(name)
	if name == "" {
		return fail("bad_name")
	}
	acc, ok := s.find(key)
	if !ok {
		return fail("not_found")
	}
	for _, a := range s.accounts() {
		if a.Key != key && strings.EqualFold(a.Name, name) {
			return fail("name_taken", a.Name)
		}
	}
	acc.Name = name
	if writeJSON(filepath.Join(s.accDir(key), "meta.json"), acc) != nil {
		return fail("internal")
	}
	return nil
}

func (s *Store) autoName() string {
	prefix := "Account"
	if s.settings().Lang == "ar" {
		prefix = "حساب"
	}
	used := map[string]bool{}
	for _, a := range s.accounts() {
		used[strings.ToLower(a.Name)] = true
	}
	for n := 1; ; n++ {
		name := prefix + " " + strconv.Itoa(n)
		if !used[strings.ToLower(name)] {
			return name
		}
	}
}

func (s *Store) AutoCapture() {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.ini)
	id := s.p.ReadAccountID()
	if err != nil || !hasLogin(data) || id == "" {
		if st := s.loadState(); st.Ignore != "" {
			st.Ignore = ""
			s.saveState(st)
		}
		return
	}
	if st := s.loadState(); st.Ignore != "" && st.Ignore != id {
		st.Ignore = ""
		s.saveState(st)
	}
	if s.refreshLocked() != "" {
		return
	}
	st := s.loadState()
	if s.settings().AutoSaveOff || st.Ignore == id {
		return
	}
	a := Account{Key: newKey(), Name: s.autoName(), ID: id, Saved: time.Now().Format(time.RFC3339)}
	if writeFile(filepath.Join(s.accDir(a.Key), iniName), data) != nil {
		return
	}
	if writeJSON(filepath.Join(s.accDir(a.Key), "meta.json"), a) != nil {
		os.RemoveAll(s.accDir(a.Key))
		return
	}
	st.Active = a.Key
	s.saveState(st)
}
