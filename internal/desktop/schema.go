package desktop

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SchemaStatus reports whether a schema directory is usable.
type SchemaStatus struct {
	Dir     string   `json:"dir"`
	Valid   bool     `json:"valid"`
	Bundles []string `json:"bundles,omitempty"`
	Message string   `json:"message,omitempty"`
}

// ValidateSchemaDir checks that dir exists and contains at least one C*.gz schema bundle.
func ValidateSchemaDir(dir string) SchemaStatus {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return SchemaStatus{Valid: false, Message: "No schema directory configured."}
	}
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return SchemaStatus{Dir: dir, Valid: false, Message: "Schema directory does not exist."}
		}
		return SchemaStatus{Dir: dir, Valid: false, Message: err.Error()}
	}
	if !info.IsDir() {
		return SchemaStatus{Dir: dir, Valid: false, Message: "Schema path is not a directory."}
	}
	bundles, err := listSchemaBundles(dir)
	if err != nil {
		return SchemaStatus{Dir: dir, Valid: false, Message: err.Error()}
	}
	if len(bundles) == 0 {
		return SchemaStatus{
			Dir:     dir,
			Valid:   false,
			Message: "No schema bundles found. Place a C27_*.gz file in this folder.",
		}
	}
	return SchemaStatus{Dir: dir, Valid: true, Bundles: bundles}
}

func listSchemaBundles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		lower := strings.ToLower(name)
		if strings.HasPrefix(strings.ToUpper(name), "C") && strings.HasSuffix(lower, ".gz") {
			out = append(out, name)
		}
	}
	return out, nil
}

func defaultSchemaCandidates(execDir string) []string {
	var dirs []string
	if wd, err := os.Getwd(); err == nil {
		dirs = append(dirs, filepath.Join(wd, "data", "schemas"))
		dirs = append(dirs, filepath.Join(wd, "schemas"))
	}
	if execDir != "" {
		dirs = append(dirs, filepath.Join(execDir, "schemas"))
		dirs = append(dirs, filepath.Join(execDir, "data", "schemas"))
	}
	return dirs
}

func firstValidSchema(candidates ...string) string {
	for _, c := range candidates {
		if c == "" {
			continue
		}
		if ValidateSchemaDir(c).Valid {
			return c
		}
	}
	return ""
}

func friendlyParseError(err error, schemaDir string) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	lower := strings.ToLower(msg)
	switch {
	case strings.Contains(lower, "schema"):
		return fmt.Errorf("%w (schema dir: %s)", err, schemaDir)
	default:
		return err
	}
}
