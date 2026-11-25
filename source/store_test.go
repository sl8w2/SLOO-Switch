package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func code(err error) string {
	if ae, ok := err.(*apiError); ok {
		return ae.Code
	}
	if err == nil {
		return ""
	}
	return "other"
}

func TestSaveSwitchRoundTripAndLaunch(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("idA", "tokA")
	a, err := s.Save("Main")
	if err != nil {
		t.Fatal(err)
	}
	e.login("idB", "tokB")
	b, err := s.Save("Alt")
	if err != nil {
		t.Fatal(err)
	}
	res, err := s.Switch(a.Key)
	if err != nil || !res.Launched || res.Name != "Main" {
		t.Fatalf("switch A: %+v %v", res, err)
	}
	if !strings.Contains(e.ini(), "tokA") || e.fp.id != "idA" {
		t.Fatalf("A not restored: %q %q", e.ini(), e.fp.id)
	}
	if len(e.fp.launched) != 1 || e.fp.launched[0] != e.exe {
		t.Fatalf("epic not launched: %v", e.fp.launched)
	}
	if _, err := s.Switch(b.Key); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(e.ini(), "tokB") || e.fp.id != "idB" {
		t.Fatalf("B not restored")
	}
}

func TestTokenStaysFreshAcrossSwitches(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("idA", "A1")
	a, _ := s.Save("A")
	e.login("idB", "B1")
	b, _ := s.Save("B")
	s.Switch(a.Key)
	e.login("idA", "A2")
	s.Switch(b.Key)
	e.login("idB", "B2")
	s.Switch(a.Key)
	if !strings.Contains(e.ini(), "A2") {
		t.Fatalf("A token went stale: %q", e.ini())
	}
	e.login("idA", "A3")
	s.Refresh()
	snap, _ := os.ReadFile(filepath.Join(s.accDir(a.Key), iniName))
	if !strings.Contains(string(snap), "A3") {
		t.Fatalf("refresh did not update snapshot")
	}
	snapB, _ := os.ReadFile(filepath.Join(s.accDir(b.Key), iniName))
	if !strings.Contains(string(snapB), "B2") {
		t.Fatalf("B snapshot should hold latest token: %q", snapB)
	}
}

func TestUnsavedAccountNeverOverwritesSaved(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("idA", "A1")
	a, _ := s.Save("A")
	e.login("idZ", "Z1")
	s.Refresh()
	snap, _ := os.ReadFile(filepath.Join(s.accDir(a.Key), iniName))
	if !strings.Contains(string(snap), "A1") {
		t.Fatalf("saved account overwritten by unrelated login")
	}
	v := s.View()
	if !v.Current.LoggedIn || v.Current.Matched != "" {
		t.Fatalf("view should report unsaved current account: %+v", v.Current)
	}
}

func TestSaveErrors(t *testing.T) {
	e := setup(t)
	s := e.app.store
	if _, err := s.Save("X"); code(err) != "no_login" {
		t.Fatalf("want no_login, got %v", err)
	}
	e.login("", "t")
	if _, err := s.Save("X"); code(err) != "no_id" {
		t.Fatalf("want no_id, got %v", err)
	}
	e.login("idA", "t")
	if _, err := s.Save("   "); code(err) != "bad_name" {
		t.Fatalf("want bad_name")
	}
	if _, err := s.Save(strings.Repeat("a", 41)); code(err) != "bad_name" {
		t.Fatalf("want bad_name for long")
	}
	if _, err := s.Save("Main"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Save("main"); code(err) != "name_taken" {
		t.Fatalf("want name_taken")
	}
	if _, err := s.Save("Other"); code(err) != "id_taken" {
		t.Fatalf("want id_taken")
	}
}

func TestUnicodeNameAndKeyValidation(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("idA", "t")
	a, err := s.Save("حسابي الأول")
	if err != nil || a.Name != "حسابي الأول" {
		t.Fatalf("unicode name failed: %v", err)
	}
	for _, k := range []string{"../x", "..", "", "ABCDEF0123456789", "zzzzzzzzzzzzzzzz"} {
		if _, err := s.Switch(k); code(err) != "not_found" {
			t.Fatalf("key %q accepted", k)
		}
		if code(s.Delete(k)) != "not_found" {
			t.Fatalf("delete key %q accepted", k)
		}
	}
	if err := s.Delete(a.Key); err != nil {
		t.Fatal(err)
	}
	if len(s.accounts()) != 0 {
		t.Fatal("not deleted")
	}
}

func TestSwitchClearsCachesBacksUpAndPrunes(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("idA", "A1")
	a, _ := s.Save("A")
	cache := s.caches[0]
	os.MkdirAll(cache, 0o700)
	os.WriteFile(filepath.Join(cache, "c"), []byte("x"), 0o600)
	for i := 0; i < 8; i++ {
		if _, err := s.Switch(a.Key); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(cache); err == nil {
		t.Fatal("cache not cleared")
	}
	entries, _ := os.ReadDir(filepath.Join(s.dir, "backups"))
	if len(entries) != keepBackups {
		t.Fatalf("backups=%d want %d", len(entries), keepBackups)
	}
}

func TestKillFailureLeavesFilesUntouched(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("idA", "A1")
	a, _ := s.Save("A")
	e.login("idB", "B1")
	e.fp.killErr = os.ErrPermission
	if _, err := s.Switch(a.Key); code(err) != "kill_failed" {
		t.Fatalf("want kill_failed, got %v", err)
	}
	if !strings.Contains(e.ini(), "B1") {
		t.Fatal("files changed despite kill failure")
	}
}

func TestSwitchWithoutEpicInstalled(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("idA", "A1")
	a, _ := s.Save("A")
	s.SetExe("")
	res, err := s.Switch(a.Key)
	if err != nil || res.Launched {
		t.Fatalf("expected switch without launch: %+v %v", res, err)
	}
	if s.View().ExeFound {
		t.Fatal("exe should not be found")
	}
	if code(s.SetExe(filepath.Join(e.local, "nope.exe"))) != "bad_path" {
		t.Fatal("bad path accepted")
	}
}

func TestActiveFlagFollowsRealLogin(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("idA", "A1")
	s.Save("A")
	e.login("idB", "B1")
	b, _ := s.Save("B")
	v := s.View()
	for _, a := range v.Accounts {
		if a.Active != (a.Key == b.Key) {
			t.Fatalf("wrong active flag: %+v", v.Accounts)
		}
	}
}

func TestNewLoginSavedAccountKeepsLatestTokenAndShowsLoginScreen(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("idA", "A1")
	a, _ := s.Save("A")
	e.login("idA", "A2")
	cache := s.caches[0]
	os.MkdirAll(cache, 0o700)
	os.WriteFile(filepath.Join(cache, "c"), []byte("x"), 0o600)
	e.fp.running = true
	res, err := s.NewLogin()
	if err != nil || !res.Launched {
		t.Fatalf("newlogin: %+v %v", res, err)
	}
	snap, _ := os.ReadFile(filepath.Join(s.accDir(a.Key), iniName))
	if !strings.Contains(string(snap), "A2") {
		t.Fatalf("latest token not captured before logout: %q", snap)
	}
	if _, err := os.Stat(s.ini); err == nil {
		t.Fatal("login file should be removed so Epic shows the login screen")
	}
	if _, err := os.Stat(cache); err == nil {
		t.Fatal("cache not cleared")
	}
	if len(e.fp.launched) != 1 {
		t.Fatalf("epic not launched: %v", e.fp.launched)
	}
	if s.View().Current.LoggedIn {
		t.Fatal("should be logged out")
	}
	e.login("idNew", "N1")
	if _, err := s.Save("New"); err != nil {
		t.Fatalf("saving the new account failed: %v", err)
	}
	snap, _ = os.ReadFile(filepath.Join(s.accDir(a.Key), iniName))
	if !strings.Contains(string(snap), "A2") {
		t.Fatal("first account snapshot damaged")
	}
	if len(s.accounts()) != 2 {
		t.Fatal("expected two accounts")
	}
}

func TestNewLoginRefusesWhenCurrentAccountUnsaved(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("idZ", "Z1")
	e.fp.running = true
	if _, err := s.NewLogin(); code(err) != "unsaved_current" {
		t.Fatalf("want unsaved_current, got %v", err)
	}
	if !e.fp.running || !strings.Contains(e.ini(), "Z1") {
		t.Fatal("must not close Epic or touch files when refusing")
	}
	if _, err := s.Save("Z"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.NewLogin(); err != nil {
		t.Fatalf("after saving it should work: %v", err)
	}
}

func TestNewLoginWhenNotLoggedInAndKillFailure(t *testing.T) {
	e := setup(t)
	s := e.app.store
	if _, err := s.NewLogin(); err != nil {
		t.Fatalf("not logged in should still open login screen: %v", err)
	}
	e.login("idA", "A1")
	s.Save("A")
	e.fp.killErr = os.ErrPermission
	if _, err := s.NewLogin(); code(err) != "kill_failed" {
		t.Fatalf("want kill_failed, got %v", err)
	}
	if !strings.Contains(e.ini(), "A1") {
		t.Fatal("files changed despite kill failure")
	}
}

func TestAutoCaptureSavesNewAccountsWithoutAnyClick(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("idA", "A1")
	s.AutoCapture()
	accs := s.accounts()
	if len(accs) != 1 || accs[0].Name != "Account 1" || accs[0].ID != "idA" {
		t.Fatalf("auto save failed: %+v", accs)
	}
	e.login("idA", "A2")
	s.AutoCapture()
	if len(s.accounts()) != 1 {
		t.Fatal("same account saved twice")
	}
	snap, _ := os.ReadFile(filepath.Join(s.accDir(accs[0].Key), iniName))
	if !strings.Contains(string(snap), "A2") {
		t.Fatalf("token not refreshed: %q", snap)
	}
	e.login("idB", "B1")
	s.AutoCapture()
	if names := []string{s.accounts()[0].Name, s.accounts()[1].Name}; names[0] != "Account 1" || names[1] != "Account 2" {
		t.Fatalf("names: %v", names)
	}
}

func TestAccountSurvivesEpicLogout(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("idA", "A1")
	s.AutoCapture()
	os.Remove(s.ini)
	s.AutoCapture()
	accs := s.accounts()
	if len(accs) != 1 {
		t.Fatal("account lost after logout")
	}
	snap, _ := os.ReadFile(filepath.Join(s.accDir(accs[0].Key), iniName))
	if !strings.Contains(string(snap), "A1") {
		t.Fatal("snapshot changed by logout")
	}
	if _, err := s.Switch(accs[0].Key); err != nil || !strings.Contains(e.ini(), "A1") {
		t.Fatalf("cannot switch back after logout: %v", err)
	}
}

func TestAutoCaptureLocalizedNameAndToggle(t *testing.T) {
	e := setup(t)
	s := e.app.store
	s.SetLang("ar")
	e.login("idA", "A1")
	s.AutoCapture()
	if s.accounts()[0].Name != "حساب 1" {
		t.Fatalf("name: %q", s.accounts()[0].Name)
	}
	if err := s.SetAutoSave(false); err != nil {
		t.Fatal(err)
	}
	if s.View().AutoSave {
		t.Fatal("autosave should be off")
	}
	e.login("idB", "B1")
	s.AutoCapture()
	if len(s.accounts()) != 1 {
		t.Fatal("captured while auto-save is off")
	}
	if s.settings().Lang != "ar" {
		t.Fatal("settings overwritten")
	}
	s.SetAutoSave(true)
	s.AutoCapture()
	if len(s.accounts()) != 2 {
		t.Fatal("not captured after enabling")
	}
}

func TestDeletedAccountIsNotResavedWhileStillLoggedIn(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("idA", "A1")
	s.AutoCapture()
	if err := s.Delete(s.accounts()[0].Key); err != nil {
		t.Fatal(err)
	}
	s.AutoCapture()
	if len(s.accounts()) != 0 {
		t.Fatal("deleted account came back")
	}
	e.login("idB", "B1")
	s.AutoCapture()
	if len(s.accounts()) != 1 {
		t.Fatal("other account should be captured")
	}
	e.login("idA", "A9")
	s.AutoCapture()
	if len(s.accounts()) != 2 {
		t.Fatal("deliberate re-login should be captured again")
	}
}

func TestAutoCaptureIgnoresMissingIdAndNewLogin(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("", "A1")
	s.AutoCapture()
	if len(s.accounts()) != 0 {
		t.Fatal("captured without id")
	}
	e.login("idA", "A1")
	s.AutoCapture()
	if _, err := s.NewLogin(); err != nil {
		t.Fatal(err)
	}
	s.AutoCapture()
	if len(s.accounts()) != 1 {
		t.Fatal("capture after new login created a ghost account")
	}
}

func TestRename(t *testing.T) {
	e := setup(t)
	s := e.app.store
	e.login("idA", "A1")
	s.AutoCapture()
	e.login("idB", "B1")
	s.AutoCapture()
	accs := s.accounts()
	if err := s.Rename(accs[0].Key, "SLOO Main"); err != nil {
		t.Fatal(err)
	}
	if code(s.Rename(accs[1].Key, "sloo main")) != "name_taken" {
		t.Fatal("duplicate name accepted")
	}
	if code(s.Rename(accs[1].Key, "  ")) != "bad_name" {
		t.Fatal("blank accepted")
	}
	if code(s.Rename("../x", "Z")) != "not_found" {
		t.Fatal("bad key accepted")
	}
	if err := s.Rename(accs[0].Key, "SLOO Main"); err != nil {
		t.Fatal("renaming to its own name must work")
	}
	found := false
	for _, a := range s.accounts() {
		if a.Name == "SLOO Main" && a.ID == "idA" {
			found = true
		}
	}
	if !found {
		t.Fatal("rename not persisted")
	}
}

func TestSettingsFieldsPreserveEachOther(t *testing.T) {
	e := setup(t)
	s := e.app.store
	s.SetAutoSave(false)
	s.SetLang("en")
	if err := s.SetExe(e.exe); err != nil {
		t.Fatal(err)
	}
	st := s.settings()
	if !st.AutoSaveOff || st.Lang != "en" || st.ExePath != e.exe {
		t.Fatalf("settings clobbered: %+v", st)
	}
}
