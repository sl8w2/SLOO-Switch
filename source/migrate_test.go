package main

import (
	"os"
	"path/filepath"
	"testing"
)

func put(t *testing.T, path, content string) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func read(path string) string {
	b, _ := os.ReadFile(path)
	return string(b)
}

func TestMigrateMovesDataAndRemovesLegacy(t *testing.T) {
	root := t.TempDir()
	legacy := filepath.Join(root, "desktop", "data")
	dst := filepath.Join(root, "pf", "data")
	put(t, filepath.Join(legacy, "accounts", "aaaaaaaaaaaaaaaa", "meta.json"), "m1")
	put(t, filepath.Join(legacy, "accounts", "aaaaaaaaaaaaaaaa", iniName), "ini1")
	put(t, filepath.Join(legacy, "state.json"), "st")
	put(t, filepath.Join(legacy, "settings.json"), "se")
	put(t, filepath.Join(legacy, "backups", "b1", iniName), "bk")
	put(t, filepath.Join(legacy, "browser", "Default", "junk"), "x")
	migrateLegacy(legacy, dst)
	if read(filepath.Join(dst, "accounts", "aaaaaaaaaaaaaaaa", iniName)) != "ini1" ||
		read(filepath.Join(dst, "state.json")) != "st" ||
		read(filepath.Join(dst, "settings.json")) != "se" ||
		read(filepath.Join(dst, "backups", "b1", iniName)) != "bk" {
		t.Fatal("data not migrated")
	}
	if _, err := os.Stat(legacy); err == nil {
		t.Fatal("legacy folder should be gone")
	}
	if _, err := os.Stat(filepath.Join(dst, "browser")); err == nil {
		t.Fatal("browser profile must not be copied")
	}
}

func TestMigrateNeverOverwritesExisting(t *testing.T) {
	root := t.TempDir()
	legacy := filepath.Join(root, "old")
	dst := filepath.Join(root, "new")
	put(t, filepath.Join(legacy, "accounts", "aaaaaaaaaaaaaaaa", iniName), "OLD")
	put(t, filepath.Join(legacy, "accounts", "bbbbbbbbbbbbbbbb", iniName), "B")
	put(t, filepath.Join(dst, "accounts", "aaaaaaaaaaaaaaaa", iniName), "NEW")
	migrateLegacy(legacy, dst)
	if read(filepath.Join(dst, "accounts", "aaaaaaaaaaaaaaaa", iniName)) != "NEW" {
		t.Fatal("existing account overwritten")
	}
	if read(filepath.Join(dst, "accounts", "bbbbbbbbbbbbbbbb", iniName)) != "B" {
		t.Fatal("missing account not copied")
	}
}

func TestMigrateLeavesForeignFoldersAndSamePath(t *testing.T) {
	root := t.TempDir()
	foreign := filepath.Join(root, "data")
	put(t, filepath.Join(foreign, "notes.txt"), "mine")
	migrateLegacy(foreign, filepath.Join(root, "dst"))
	if read(filepath.Join(foreign, "notes.txt")) != "mine" {
		t.Fatal("foreign folder touched")
	}
	same := filepath.Join(root, "same")
	put(t, filepath.Join(same, "state.json"), "s")
	migrateLegacy(same, same)
	if read(filepath.Join(same, "state.json")) != "s" {
		t.Fatal("same-path migration damaged data")
	}
}
