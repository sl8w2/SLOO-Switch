package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func serve(t *testing.T) (*env, *httptest.Server) {
	e := setup(t)
	ts := httptest.NewServer(e.app.handler())
	u, _ := url.Parse(ts.URL)
	port := 0
	for _, c := range u.Port() {
		port = port*10 + int(c-'0')
	}
	e.app.port = port
	t.Cleanup(ts.Close)
	return e, ts
}

func post(t *testing.T, ts *httptest.Server, path, token, body string, host string) (int, map[string]any) {
	req, _ := http.NewRequest("POST", ts.URL+path, bytes.NewBufferString(body))
	if token != "" {
		req.Header.Set("X-SLOO", token)
	}
	if host != "" {
		req.Host = host
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out map[string]any
	json.NewDecoder(resp.Body).Decode(&out)
	return resp.StatusCode, out
}

func TestAPIRequiresTokenAndHost(t *testing.T) {
	e, ts := serve(t)
	if c, _ := post(t, ts, "/api/state", "", "{}", ""); c != 403 {
		t.Fatalf("no token: %d", c)
	}
	if c, _ := post(t, ts, "/api/state", "wrong", "{}", ""); c != 403 {
		t.Fatalf("bad token: %d", c)
	}
	if c, _ := post(t, ts, "/api/state", e.app.token, "{}", "evil.example:80"); c != 403 {
		t.Fatalf("bad host: %d", c)
	}
	if c, out := post(t, ts, "/api/state", e.app.token, "{}", ""); c != 200 || out["version"] != version {
		t.Fatalf("good request failed: %d %v", c, out)
	}
	resp, _ := http.Get(ts.URL + "/api/state")
	if resp.StatusCode != 405 {
		t.Fatalf("GET should be rejected: %d", resp.StatusCode)
	}
}

func TestAPIFlow(t *testing.T) {
	e, ts := serve(t)
	e.login("idA", "A1")
	_, out := post(t, ts, "/api/save", e.app.token, `{"name":"Main"}`, "")
	if out["name"] != "Main" {
		t.Fatalf("save: %v", out)
	}
	_, out = post(t, ts, "/api/save", e.app.token, `{"name":"Main"}`, "")
	if out["error"] != "name_taken" {
		t.Fatalf("dup: %v", out)
	}
	_, st := post(t, ts, "/api/state", e.app.token, "{}", "")
	accs := st["accounts"].([]any)
	k := accs[0].(map[string]any)["key"].(string)
	_, out = post(t, ts, "/api/switch", e.app.token, `{"key":"`+k+`"}`, "")
	if out["launched"] != true {
		t.Fatalf("switch: %v", out)
	}
	_, out = post(t, ts, "/api/delete", e.app.token, `{"key":"../../x"}`, "")
	if out["error"] != "not_found" {
		t.Fatalf("traversal: %v", out)
	}
	_, out = post(t, ts, "/api/delete", e.app.token, `{"key":"`+k+`"}`, "")
	if out["ok"] != true {
		t.Fatalf("delete: %v", out)
	}
}

func TestOpenWhitelist(t *testing.T) {
	e, ts := serve(t)
	post(t, ts, "/api/open", e.app.token, `{"url":"https://evil.example"}`, "")
	post(t, ts, "/api/open", e.app.token, `{"url":"https://www.tiktok.com/@sl8w2"}`, "")
	post(t, ts, "/api/open", e.app.token, `{"url":"https://www.instagram.com/sl8w2"}`, "")
	if len(e.fp.opened) != 2 {
		t.Fatalf("opened=%v", e.fp.opened)
	}
}

func TestIndexAndWho(t *testing.T) {
	e, ts := serve(t)
	resp, _ := http.Get(ts.URL + "/")
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	body := buf.String()
	if !strings.Contains(body, e.app.token) || strings.Contains(body, "__TOKEN__") {
		t.Fatal("token not injected")
	}
	for _, s := range []string{"SLOO", "sl8w2", "TikTok", "Instagram", `id="theme"`, `data-theme=dark`, "prefers-reduced-motion", `role="dialog"`} {
		if !strings.Contains(body, s) {
			t.Fatalf("credit %q missing", s)
		}
	}
	if resp.Header.Get("Content-Security-Policy") == "" {
		t.Fatal("no CSP")
	}
	r2, _ := http.Get(ts.URL + "/who")
	buf.Reset()
	buf.ReadFrom(r2.Body)
	if buf.String() != "SLOO-SWITCH" {
		t.Fatal("who")
	}
	r3, _ := http.Get(ts.URL + "/nope")
	if r3.StatusCode != 404 {
		t.Fatal("404 expected")
	}
}

func TestWatchdogByeAndPingCancel(t *testing.T) {
	e, _ := serve(t)
	e.app.touch(false)
	e.app.touch(true)
	e.app.touch(false)
	e.app.mu.Lock()
	cleared := e.app.bye.IsZero()
	e.app.mu.Unlock()
	if !cleared {
		t.Fatal("ping should cancel bye")
	}
	e.app.touch(true)
	e.app.mu.Lock()
	e.app.bye = time.Now().Add(-time.Second)
	e.app.mu.Unlock()
	done := make(chan struct{})
	go func() { e.app.watch(); close(done) }()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("watchdog did not stop after bye")
	}
}

func TestIconServedAndValid(t *testing.T) {
	_, ts := serve(t)
	resp, err := http.Get(ts.URL + "/icon.png")
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	b := buf.Bytes()
	if resp.StatusCode != 200 || resp.Header.Get("Content-Type") != "image/png" {
		t.Fatalf("icon status/type: %d %s", resp.StatusCode, resp.Header.Get("Content-Type"))
	}
	if len(b) < 1000 || string(b[1:4]) != "PNG" {
		t.Fatalf("icon is not a png: %d bytes", len(b))
	}
}

func TestAPINewLogin(t *testing.T) {
	e, ts := serve(t)
	e.login("idA", "A1")
	_, out := post(t, ts, "/api/newlogin", e.app.token, "{}", "")
	if out["error"] != "unsaved_current" {
		t.Fatalf("expected unsaved_current: %v", out)
	}
	post(t, ts, "/api/save", e.app.token, `{"name":"Main"}`, "")
	_, out = post(t, ts, "/api/newlogin", e.app.token, "{}", "")
	if out["launched"] != true {
		t.Fatalf("newlogin: %v", out)
	}
}

func TestTrackedAppIgnoresMissingHeartbeats(t *testing.T) {
	e, _ := serve(t)
	e.app.setTracked(true)
	e.app.mu.Lock()
	e.app.last = time.Now().Add(-10 * time.Minute)
	e.app.bye = time.Now().Add(-time.Minute)
	e.app.mu.Unlock()
	done := make(chan struct{})
	go func() { e.app.watch(); close(done) }()
	select {
	case <-done:
		t.Fatal("tracked app exited without its window closing")
	case <-time.After(2500 * time.Millisecond):
	}
	e.app.setTracked(false)
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("untracked app should fall back to heartbeat exit")
	}
}

func TestAPIRenameAutoSaveLang(t *testing.T) {
	e, ts := serve(t)
	e.login("idA", "A1")
	e.app.store.AutoCapture()
	_, st := post(t, ts, "/api/state", e.app.token, "{}", "")
	if st["autoSave"] != true {
		t.Fatalf("autoSave default: %v", st["autoSave"])
	}
	k := st["accounts"].([]any)[0].(map[string]any)["key"].(string)
	_, out := post(t, ts, "/api/rename", e.app.token, `{"key":"`+k+`","name":"Mine"}`, "")
	if out["ok"] != true {
		t.Fatalf("rename: %v", out)
	}
	post(t, ts, "/api/autosave", e.app.token, `{"on":false}`, "")
	post(t, ts, "/api/lang", e.app.token, `{"lang":"ar"}`, "")
	_, st = post(t, ts, "/api/state", e.app.token, "{}", "")
	if st["autoSave"] != false {
		t.Fatalf("autoSave should be off: %v", st["autoSave"])
	}
	if e.app.store.settings().Lang != "ar" {
		t.Fatal("lang not stored")
	}
}
