package main

import (
	"os"
	"path/filepath"
	"testing"
)

type fakePlatform struct {
	id       string
	running  bool
	killErr  error
	launched []string
	opened   []string
	dataDir  string
}

func newPlatform() platform { return &fakePlatform{} }

func (f *fakePlatform) ReadAccountID() string          { return f.id }
func (f *fakePlatform) WriteAccountID(id string) error { f.id = id; return nil }
func (f *fakePlatform) EpicRunning() bool              { return f.running }
func (f *fakePlatform) KillEpic() error {
	if f.killErr != nil {
		return f.killErr
	}
	f.running = false
	return nil
}
func (f *fakePlatform) Launch(exe string) error {
	f.launched = append(f.launched, exe)
	f.running = true
	return nil
}
func (f *fakePlatform) OpenUI(url, profile string) <-chan struct{} { return nil }
func (f *fakePlatform) OpenURL(url string)                         { f.opened = append(f.opened, url) }
func (f *fakePlatform) DataDir(exeDir string) string               { return f.dataDir }
func (f *fakePlatform) ElevatedSetup(sid string) int               { return 0 }

type env struct {
	t     *testing.T
	fp    *fakePlatform
	app   *App
	local string
	exe   string
}

func setup(t *testing.T) *env {
	root := t.TempDir()
	local := filepath.Join(root, "local")
	t.Setenv("LOCALAPPDATA", local)
	t.Setenv("ProgramFiles", "")
	t.Setenv("ProgramFiles(x86)", "")
	fp := &fakePlatform{}
	app := newApp(filepath.Join(root, "data"), fp)
	exe := filepath.Join(root, "EpicGamesLauncher.exe")
	os.WriteFile(exe, []byte("x"), 0o600)
	if err := app.store.SetExe(exe); err != nil {
		t.Fatal(err)
	}
	return &env{t: t, fp: fp, app: app, local: local, exe: exe}
}

func (e *env) login(id, token string) {
	e.fp.id = id
	os.MkdirAll(filepath.Dir(e.app.store.ini), 0o700)
	os.WriteFile(e.app.store.ini, []byte("[RememberMe]\nEnable=True\nData=\""+token+"\"\n"), 0o600)
}

func (e *env) ini() string {
	b, _ := os.ReadFile(e.app.store.ini)
	return string(b)
}
