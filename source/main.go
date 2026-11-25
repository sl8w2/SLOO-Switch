package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const fixedPort = 47821

func alreadyRunning() bool {
	c := http.Client{Timeout: 800 * time.Millisecond}
	resp, err := c.Get(fmt.Sprintf("http://127.0.0.1:%d/who", fixedPort))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 32))
	return string(b) == "SLOO-SWITCH"
}

func main() {
	p := newPlatform()
	if len(os.Args) == 3 && os.Args[1] == "--setup" {
		os.Exit(p.ElevatedSetup(os.Args[2]))
	}
	exe, err := os.Executable()
	if err != nil {
		return
	}
	exeDir := filepath.Dir(exe)
	data := p.DataDir(exeDir)
	migrateLegacy(filepath.Join(exeDir, "data"), data)
	profile := filepath.Join(data, "browser")
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", fixedPort))
	if err != nil {
		if alreadyRunning() {
			p.OpenUI(fmt.Sprintf("http://127.0.0.1:%d/", fixedPort), profile)
			return
		}
		if ln, err = net.Listen("tcp", "127.0.0.1:0"); err != nil {
			return
		}
	}
	app := newApp(data, p)
	app.port = ln.Addr().(*net.TCPAddr).Port
	srv := &http.Server{Handler: app.handler(), ReadHeaderTimeout: 5 * time.Second}
	go srv.Serve(ln)
	app.store.Refresh()
	go app.watch()
	done := p.OpenUI(fmt.Sprintf("http://127.0.0.1:%d/", app.port), profile)
	if done != nil {
		app.setTracked(true)
		started := time.Now()
		go func() {
			<-done
			if time.Since(started) > 5*time.Second {
				app.stop()
			} else {
				app.setTracked(false)
			}
		}()
	}
	<-app.quit
	srv.Close()
}
