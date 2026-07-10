package dynasty

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

var (
	sharedTestSaveOnce sync.Once
	sharedTestSaveFile *File
	sharedTestSaveErr  error
)

// openTestSave returns a parsed test dynasty save, opened once and shared across
// the package. Integration tests must not mutate table rows in ways that affect
// other tests.
func openTestSave(t *testing.T) *File {
	t.Helper()
	sharedTestSaveOnce.Do(func() {
		savePath := filepath.Join("..", "data", "DYNASTY-TESTSAVE-27")
		schemaDir := filepath.Join("..", "data")
		if _, err := os.Stat(savePath); err != nil {
			sharedTestSaveErr = err
			return
		}
		settings := DefaultSettings()
		settings.SchemaDir = schemaDir
		settings.AutoParse = true
		sharedTestSaveFile, sharedTestSaveErr = Open(savePath, &settings)
	})
	if sharedTestSaveErr != nil {
		if os.IsNotExist(sharedTestSaveErr) {
			t.Skip("test save not available:", filepath.Join("..", "data", "DYNASTY-TESTSAVE-27"))
		}
		t.Fatal(sharedTestSaveErr)
	}
	return sharedTestSaveFile
}

var (
	sharedSeasonSaveOnce sync.Once
	sharedSeasonSaveFile *File
	sharedSeasonSaveErr  error
)

// openSeasonSave returns a parsed full-season test save (DYNASTY-2026OFFLINEFINAL),
// opened once and shared across the package.
func openSeasonSave(t *testing.T) *File {
	t.Helper()
	sharedSeasonSaveOnce.Do(func() {
		savePath := filepath.Join("..", "data", "DYNASTY-2026OFFLINEFINAL")
		schemaDir := filepath.Join("..", "data")
		if _, err := os.Stat(savePath); err != nil {
			sharedSeasonSaveErr = err
			return
		}
		settings := DefaultSettings()
		settings.SchemaDir = schemaDir
		settings.AutoParse = true
		sharedSeasonSaveFile, sharedSeasonSaveErr = Open(savePath, &settings)
	})
	if sharedSeasonSaveErr != nil {
		if os.IsNotExist(sharedSeasonSaveErr) {
			t.Skip("season save not available:", filepath.Join("..", "data", "DYNASTY-2026OFFLINEFINAL"))
		}
		t.Fatal(sharedSeasonSaveErr)
	}
	return sharedSeasonSaveFile
}

// skipIfShortIntegration skips heavy tests that parse a real save and export
// large sections. Run with -short for fast unit tests; omit -short for integration.
func skipIfShortIntegration(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test (omit -short to run)")
	}
}
