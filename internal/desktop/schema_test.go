package desktop_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leaguelines/cfb-dynasty/internal/desktop"
	"github.com/leaguelines/cfb-dynasty/internal/desktop/saves"
)

func TestValidateSchemaDir(t *testing.T) {
	dir := t.TempDir()
	st := desktop.ValidateSchemaDir(dir)
	if st.Valid {
		t.Fatal("empty dir should be invalid")
	}

	if err := os.WriteFile(filepath.Join(dir, "C27_472_0.gz"), []byte("not-a-real-bundle"), 0o644); err != nil {
		t.Fatal(err)
	}
	st = desktop.ValidateSchemaDir(dir)
	if !st.Valid {
		t.Fatalf("expected valid schema dir, got %#v", st)
	}
	if len(st.Bundles) != 1 || st.Bundles[0] != "C27_472_0.gz" {
		t.Fatalf("bundles = %#v", st.Bundles)
	}
}

func TestLooksLikeSaveViaDiscoverEmpty(t *testing.T) {
	// Discover should tolerate missing default dirs.
	_, err := saves.Discover()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDefaultWindowsSavesDir(t *testing.T) {
	dir := saves.DefaultWindowsSavesDir()
	if dir == "" {
		t.Fatal("expected non-empty default saves dir")
	}
	if filepath.Base(dir) != "saves" {
		t.Fatalf("unexpected dir %q", dir)
	}
}
