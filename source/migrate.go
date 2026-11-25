package main

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var legacyItems = []string{"accounts", "backups", "state.json", "settings.json"}

func isOurData(dir string) bool {
	for _, n := range []string{"accounts", "state.json"} {
		if _, err := os.Stat(filepath.Join(dir, n)); err == nil {
			return true
		}
	}
	return false
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

func copyMissing(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o700)
		}
		if _, err := os.Stat(target); err == nil {
			return nil
		}
		return copyFile(p, target)
	})
}

func migrateLegacy(legacy, dst string) {
	if filepath.Clean(legacy) == filepath.Clean(dst) || !isOurData(legacy) {
		return
	}
	for _, n := range legacyItems {
		src := filepath.Join(legacy, n)
		fi, err := os.Stat(src)
		if err != nil {
			continue
		}
		target := filepath.Join(dst, n)
		if fi.IsDir() {
			err = copyMissing(src, target)
		} else if _, statErr := os.Stat(target); statErr != nil {
			err = copyFile(src, target)
		}
		if err != nil {
			return
		}
	}
	for _, n := range append(append([]string{}, legacyItems...), "browser") {
		os.RemoveAll(filepath.Join(legacy, n))
	}
	os.Remove(legacy)
}
