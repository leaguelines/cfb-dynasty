// Package prefs persists desktop GUI settings (schema directory, etc.).
package prefs

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const appDirName = "cfb-dynasty-gui"

// Prefs holds user preferences for the desktop app.
type Prefs struct {
	SchemaDir string `json:"schemaDir,omitempty"`
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, appDirName)
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "prefs.json"), nil
}

// Load reads prefs from the user config directory. Missing file returns empty Prefs.
func Load() (Prefs, error) {
	path, err := configPath()
	if err != nil {
		return Prefs{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Prefs{}, nil
		}
		return Prefs{}, err
	}
	var p Prefs
	if err := json.Unmarshal(data, &p); err != nil {
		return Prefs{}, err
	}
	return p, nil
}

// Save writes prefs to the user config directory.
func Save(p Prefs) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
