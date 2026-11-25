package main

type platform interface {
	ReadAccountID() string
	WriteAccountID(id string) error
	EpicRunning() bool
	KillEpic() error
	Launch(exe string) error
	OpenUI(url, profileDir string) <-chan struct{}
	OpenURL(url string)
	DataDir(exeDir string) string
	ElevatedSetup(sid string) int
}
