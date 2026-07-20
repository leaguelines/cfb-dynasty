// Package saves discovers dynasty save files on disk.
package saves

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// SaveFile is a discovered dynasty save candidate.
type SaveFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
}

// DefaultWindowsSavesDir returns the EA CFB 27 default saves directory on Windows.
func DefaultWindowsSavesDir() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, "Documents", "EA SPORTS CFB27", "saves")
}

// CandidateDirs returns directories to scan for saves, in priority order.
func CandidateDirs() []string {
	var dirs []string
	seen := map[string]bool{}
	add := func(d string) {
		if d == "" || seen[d] {
			return
		}
		seen[d] = true
		dirs = append(dirs, d)
	}

	if runtime.GOOS == "windows" {
		add(DefaultWindowsSavesDir())
	} else {
		// Native Windows path when running on Windows; elsewhere try common Wine layouts.
		add(DefaultWindowsSavesDir())
		home, _ := os.UserHomeDir()
		if home != "" {
			add(filepath.Join(home, ".wine", "drive_c", "users", os.Getenv("USER"), "Documents", "EA SPORTS CFB27", "saves"))
			add(filepath.Join(home, ".wine", "drive_c", "users", "steamuser", "Documents", "EA SPORTS CFB27", "saves"))
		}
	}
	return dirs
}

// Discover lists plausible dynasty save files under candidate directories.
func Discover() ([]SaveFile, error) {
	var out []SaveFile
	for _, dir := range CandidateDirs() {
		found, err := scanDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		out = append(out, found...)
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

func scanDir(dir string) ([]SaveFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []SaveFile
	for _, e := range entries {
		if e.IsDir() {
			// One level of nesting (some installs use subfolders).
			sub := filepath.Join(dir, e.Name())
			subEntries, err := os.ReadDir(sub)
			if err != nil {
				continue
			}
			for _, se := range subEntries {
				if se.IsDir() {
					continue
				}
				if sf, ok := asSave(filepath.Join(sub, se.Name()), se); ok {
					out = append(out, sf)
				}
			}
			continue
		}
		if sf, ok := asSave(filepath.Join(dir, e.Name()), e); ok {
			out = append(out, sf)
		}
	}
	return out, nil
}

func asSave(path string, e os.DirEntry) (SaveFile, bool) {
	name := e.Name()
	if !looksLikeSave(name) {
		return SaveFile{}, false
	}
	info, err := e.Info()
	if err != nil {
		return SaveFile{}, false
	}
	// Ignore tiny junk files; dynasty saves are typically multi-MB.
	if info.Size() < 64*1024 {
		return SaveFile{}, false
	}
	return SaveFile{Name: name, Path: path, Size: info.Size()}, true
}

func looksLikeSave(name string) bool {
	lower := strings.ToLower(name)
	switch {
	case strings.HasSuffix(lower, ".sav"):
		return true
	case strings.HasSuffix(lower, ".json"), strings.HasSuffix(lower, ".csv"),
		strings.HasSuffix(lower, ".txt"), strings.HasSuffix(lower, ".png"),
		strings.HasSuffix(lower, ".jpg"), strings.HasSuffix(lower, ".log"),
		strings.HasSuffix(lower, ".gz"), strings.HasSuffix(lower, ".zip"):
		return false
	case strings.HasPrefix(name, "."):
		return false
	default:
		// EA CFB dynasty saves are often extensionless.
		return !strings.Contains(name, ".")
	}
}
