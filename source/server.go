package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

var openAllowed = map[string]bool{
	"https://www.tiktok.com/@sl8w2":   true,
	"https://www.instagram.com/sl8w2": true,
}

type App struct {
	store *Store
	p     platform
	token string
	port  int
	start time.Time
	quit  chan struct{}
	once  sync.Once
	mu    sync.Mutex
	last  time.Time
	bye   time.Time
	track bool
}

func newApp(dir string, p platform) *App {
	b := make([]byte, 24)
	rand.Read(b)
	return &App{
		store: newStore(dir, p),
		p:     p,
		token: hex.EncodeToString(b),
		start: time.Now(),
		quit:  make(chan struct{}),
	}
}

func (a *App) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.index)
	mux.HandleFunc("/icon.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "max-age=86400")
		w.Write(iconPNG)
	})
	mux.HandleFunc("/who", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("SLOO-SWITCH"))
	})
	mux.HandleFunc("/api/", a.api)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != fmt.Sprintf("127.0.0.1:%d", a.port) && r.Host != fmt.Sprintf("localhost:%d", a.port) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		mux.ServeHTTP(w, r)
	})
}

func (a *App) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	page := strings.NewReplacer("__TOKEN__", a.token, "__VERSION__", version).Replace(pageHTML)
	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("Cache-Control", "no-store")
	h.Set("X-Frame-Options", "DENY")
	h.Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; script-src 'unsafe-inline'; connect-src 'self'; img-src 'self'")
	w.Write([]byte(page))
}

func (a *App) setTracked(v bool) {
	a.mu.Lock()
	a.track = v
	a.mu.Unlock()
}

func (a *App) touch(bye bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if bye {
		a.bye = time.Now().Add(5 * time.Second)
		return
	}
	a.last = time.Now()
	a.bye = time.Time{}
}

func reply(w http.ResponseWriter, v any, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		out := map[string]string{"error": "internal"}
		if ae, ok := err.(*apiError); ok {
			out["error"] = ae.Code
			out["arg"] = ae.Arg
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	json.NewEncoder(w).Encode(v)
}

func (a *App) api(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	tok := r.Header.Get("X-SLOO")
	if tok == "" {
		tok = r.URL.Query().Get("t")
	}
	if subtle.ConstantTimeCompare([]byte(tok), []byte(a.token)) != 1 {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var in struct {
		Name string `json:"name"`
		Key  string `json:"key"`
		Path string `json:"path"`
		On   bool   `json:"on"`
		Lang string `json:"lang"`
		URL  string `json:"url"`
	}
	json.NewDecoder(r.Body).Decode(&in)
	ok := map[string]bool{"ok": true}
	switch strings.TrimPrefix(r.URL.Path, "/api/") {
	case "ping":
		a.touch(false)
		reply(w, ok, nil)
	case "bye":
		a.touch(true)
		reply(w, ok, nil)
	case "state":
		reply(w, a.store.View(), nil)
	case "save":
		acc, err := a.store.Save(in.Name)
		reply(w, acc, err)
	case "switch":
		res, err := a.store.Switch(in.Key)
		reply(w, res, err)
	case "newlogin":
		res, err := a.store.NewLogin()
		reply(w, res, err)
	case "delete":
		reply(w, ok, a.store.Delete(in.Key))
	case "rename":
		reply(w, ok, a.store.Rename(in.Key, in.Name))
	case "autosave":
		reply(w, ok, a.store.SetAutoSave(in.On))
	case "lang":
		a.store.SetLang(in.Lang)
		reply(w, ok, nil)
	case "exe":
		reply(w, ok, a.store.SetExe(in.Path))
	case "open":
		if openAllowed[in.URL] {
			a.p.OpenURL(in.URL)
		}
		reply(w, ok, nil)
	default:
		http.NotFound(w, r)
	}
}

func (a *App) stop() {
	a.store.Refresh()
	a.once.Do(func() { close(a.quit) })
}

func (a *App) watch() {
	tick := time.NewTicker(time.Second)
	refresh := time.NewTicker(4 * time.Second)
	defer tick.Stop()
	defer refresh.Stop()
	for {
		select {
		case <-tick.C:
			a.mu.Lock()
			now := time.Now()
			done := !a.track && ((!a.bye.IsZero() && now.After(a.bye)) ||
				(a.last.IsZero() && now.Sub(a.start) > 90*time.Second) ||
				(!a.last.IsZero() && now.Sub(a.last) > 150*time.Second))
			a.mu.Unlock()
			if done {
				a.stop()
				return
			}
		case <-refresh.C:
			a.store.AutoCapture()
		case <-a.quit:
			return
		}
	}
}
