package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	idKey  = `Software\Epic Games\Unreal Engine\Identifiers`
	idName = "AccountId"
)

var (
	advapi       = syscall.NewLazyDLL("advapi32.dll")
	regCreateKey = advapi.NewProc("RegCreateKeyExW")
	regSetValue  = advapi.NewProc("RegSetValueExW")
	shell        = syscall.NewLazyDLL("shell32.dll")
	shellExec    = shell.NewProc("ShellExecuteW")
	epicNames    = map[string]bool{
		"epicgameslauncher.exe":                 true,
		"epicwebhelper.exe":                     true,
		"eosoverlayrenderer-win64-shipping.exe": true,
		"eosoverlayrenderer-win32-shipping.exe": true,
	}
)

type winPlatform struct{}

func newPlatform() platform { return winPlatform{} }

func (winPlatform) ReadAccountID() string {
	k, _ := syscall.UTF16PtrFromString(idKey)
	var h syscall.Handle
	if syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_USER, k, 0, syscall.KEY_READ, &h) != nil {
		return ""
	}
	defer syscall.RegCloseKey(h)
	n, _ := syscall.UTF16PtrFromString(idName)
	var typ, size uint32
	if syscall.RegQueryValueEx(h, n, nil, &typ, nil, &size) != nil || size == 0 {
		return ""
	}
	buf := make([]uint16, size/2+1)
	size = uint32(len(buf) * 2)
	if syscall.RegQueryValueEx(h, n, nil, &typ, (*byte)(unsafe.Pointer(&buf[0])), &size) != nil {
		return ""
	}
	return strings.TrimSpace(syscall.UTF16ToString(buf))
}

func (winPlatform) WriteAccountID(id string) error {
	k, _ := syscall.UTF16PtrFromString(idKey)
	var h syscall.Handle
	var disp uint32
	r, _, _ := regCreateKey.Call(uintptr(syscall.HKEY_CURRENT_USER), uintptr(unsafe.Pointer(k)), 0, 0, 0, 0x20006, 0, uintptr(unsafe.Pointer(&h)), uintptr(unsafe.Pointer(&disp)))
	if r != 0 {
		return syscall.Errno(r)
	}
	defer syscall.RegCloseKey(h)
	n, _ := syscall.UTF16PtrFromString(idName)
	v, err := syscall.UTF16FromString(id)
	if err != nil {
		return err
	}
	r, _, _ = regSetValue.Call(uintptr(h), uintptr(unsafe.Pointer(n)), 0, 1, uintptr(unsafe.Pointer(&v[0])), uintptr(len(v)*2))
	if r != 0 {
		return syscall.Errno(r)
	}
	return nil
}

func epicPids() []uint32 {
	snap, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil
	}
	defer syscall.CloseHandle(snap)
	var pe syscall.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	if syscall.Process32First(snap, &pe) != nil {
		return nil
	}
	var out []uint32
	for {
		if epicNames[strings.ToLower(syscall.UTF16ToString(pe.ExeFile[:]))] {
			out = append(out, pe.ProcessID)
		}
		if syscall.Process32Next(snap, &pe) != nil {
			break
		}
	}
	return out
}

func (winPlatform) EpicRunning() bool { return len(epicPids()) > 0 }

func (winPlatform) KillEpic() error {
	for i := 0; i < 6; i++ {
		pids := epicPids()
		if len(pids) == 0 {
			return nil
		}
		for _, pid := range pids {
			h, err := syscall.OpenProcess(syscall.PROCESS_TERMINATE|syscall.SYNCHRONIZE, false, pid)
			if err != nil {
				continue
			}
			syscall.TerminateProcess(h, 1)
			syscall.WaitForSingleObject(h, 3000)
			syscall.CloseHandle(h)
		}
		time.Sleep(300 * time.Millisecond)
	}
	if len(epicPids()) > 0 {
		return syscall.EBUSY
	}
	return nil
}

func (winPlatform) Launch(exe string) error {
	c := exec.Command(exe)
	c.Dir = filepath.Dir(exe)
	c.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x00000208}
	if err := c.Start(); err != nil {
		return err
	}
	return c.Process.Release()
}

func browserPaths() []string {
	rel := []string{
		`Google\Chrome\Application\chrome.exe`,
		`BraveSoftware\Brave-Browser\Application\brave.exe`,
		`Vivaldi\Application\vivaldi.exe`,
		`Chromium\Application\chrome.exe`,
	}
	var out []string
	for _, env := range []string{"ProgramFiles", "ProgramFiles(x86)", "LOCALAPPDATA"} {
		base := os.Getenv(env)
		if base == "" {
			continue
		}
		for _, r := range rel {
			out = append(out, filepath.Join(base, r))
		}
	}
	return out
}

func (w winPlatform) OpenUI(url, profile string) <-chan struct{} {
	for _, b := range browserPaths() {
		if _, err := os.Stat(b); err != nil {
			continue
		}
		c := exec.Command(b, "--app="+url, "--window-size=1120,780", "--user-data-dir="+profile, "--no-first-run", "--no-default-browser-check", "--disable-extensions", "--disable-background-mode", "--disable-background-networking", "--disable-sync", "--disable-component-update", "--disable-default-apps", "--disable-breakpad", "--disable-domain-reliability", "--metrics-recording-only", "--disk-cache-size=1", "--media-cache-size=1", "--disable-features=Translate,MediaRouter,OptimizationHints")
		if c.Start() != nil {
			continue
		}
		done := make(chan struct{})
		go func() {
			c.Wait()
			close(done)
		}()
		return done
	}
	w.OpenURL(url)
	return nil
}

func (winPlatform) OpenURL(url string) {
	verb, _ := syscall.UTF16PtrFromString("open")
	target, _ := syscall.UTF16PtrFromString(url)
	shellExec.Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(target)), 0, 0, 1)
}

var sidRe = regexp.MustCompile(`^S-1-[0-9]+(-[0-9]+)+$`)

func writableDir(dir string) bool {
	if os.MkdirAll(dir, 0o700) != nil {
		return false
	}
	f, err := os.CreateTemp(dir, ".w")
	if err != nil {
		return false
	}
	name := f.Name()
	f.Close()
	os.Remove(name)
	return true
}

func hasData(dir string) bool {
	for _, n := range []string{"state.json", "accounts"} {
		if _, err := os.Stat(filepath.Join(dir, n)); err == nil {
			return true
		}
	}
	return false
}

func currentSID() string {
	t, err := syscall.OpenCurrentProcessToken()
	if err != nil {
		return ""
	}
	defer t.Close()
	u, err := t.GetTokenUser()
	if err != nil {
		return ""
	}
	s, err := u.User.Sid.String()
	if err != nil {
		return ""
	}
	return s
}

func requestSetup(target string) bool {
	sid := currentSID()
	exe, err := os.Executable()
	if sid == "" || err != nil {
		return false
	}
	verb, _ := syscall.UTF16PtrFromString("runas")
	file, _ := syscall.UTF16PtrFromString(exe)
	args, _ := syscall.UTF16PtrFromString("--setup " + sid)
	r, _, _ := shellExec.Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(file)), uintptr(unsafe.Pointer(args)), 0, 0)
	if r <= 32 {
		return false
	}
	end := time.Now().Add(120 * time.Second)
	for time.Now().Before(end) {
		if writableDir(target) {
			return true
		}
		time.Sleep(300 * time.Millisecond)
	}
	return false
}

func (winPlatform) DataDir(exeDir string) string {
	local := filepath.Join(os.Getenv("LOCALAPPDATA"), "SLOO Switch", "data")
	pf := os.Getenv("ProgramFiles")
	if pf == "" {
		return local
	}
	target := filepath.Join(pf, "SLOO Switch", "data")
	if writableDir(target) {
		return target
	}
	if hasData(local) {
		return local
	}
	if requestSetup(target) {
		return target
	}
	return local
}

func (winPlatform) ElevatedSetup(sid string) int {
	pf := os.Getenv("ProgramFiles")
	if pf == "" || !sidRe.MatchString(sid) {
		return 2
	}
	root := filepath.Join(pf, "SLOO Switch")
	if os.MkdirAll(filepath.Join(root, "data"), 0o700) != nil {
		return 3
	}
	c := exec.Command(filepath.Join(os.Getenv("SystemRoot"), "System32", "icacls.exe"), root,
		"/inheritance:r", "/grant:r", "*"+sid+":(OI)(CI)F", "*S-1-5-18:(OI)(CI)F", "*S-1-5-32-544:(OI)(CI)F", "/T", "/C", "/Q")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if c.Run() != nil {
		return 4
	}
	return 0
}
