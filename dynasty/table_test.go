package dynasty

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseMinimalSave(t *testing.T) {
	path := t.TempDir() + "/test.sav"
	payload := append([]byte("FrTk"), make([]byte, 60)...)
	if err := writeFile(path, payload); err != nil {
		t.Fatal(err)
	}

	file, err := Open(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Parse(); err != nil {
		t.Fatalf("Parse() = %v", err)
	}
	if !file.Loaded() {
		t.Fatal("expected loaded file")
	}
	if len(file.Tables()) != 0 {
		t.Fatalf("tables = %d, want 0", len(file.Tables()))
	}
}

func TestDiscoverTablesRealSave(t *testing.T) {
	savePath := filepath.Join("..", "data", "DYNASTY-TESTSAVE-27")
	schemaDir := filepath.Join("..", "data")
	if _, err := os.Stat(savePath); err != nil {
		t.Skip("test save not available:", savePath)
	}

	settings := DefaultSettings()
	settings.SchemaDir = schemaDir
	file, err := Open(savePath, &settings)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Parse(); err != nil {
		t.Fatal(err)
	}

	tables := file.Tables()
	if len(tables) < 2000 {
		t.Fatalf("tables = %d, want at least 2000", len(tables))
	}

	seasonGame, ok := file.GetTableByName("SeasonGame")
	if !ok {
		t.Fatal("SeasonGame table not found")
	}
	// Prefer the main schedule table with the most allocated rows.
	for _, table := range tables {
		if table.Name() != "SeasonGame" {
			continue
		}
		if table.AllocatedRecordCount() > seasonGame.AllocatedRecordCount() {
			seasonGame = &table
		}
	}

	if seasonGame.AllocatedRecordCount() != 983 {
		t.Fatalf("SeasonGame records = %d, want 983", seasonGame.AllocatedRecordCount())
	}
	if seasonGame.ActiveRecordCount() != 934 {
		t.Fatalf("SeasonGame active = %d, want 934", seasonGame.ActiveRecordCount())
	}
	if seasonGame.Header.NumMembers != 69 {
		t.Fatalf("SeasonGame members = %d, want 69", seasonGame.Header.NumMembers)
	}
	if seasonGame.Header.RecordSize != 100 {
		t.Fatalf("SeasonGame record size = %d, want 100", seasonGame.Header.RecordSize)
	}
	if seasonGame.Schema == nil {
		t.Fatal("SeasonGame schema not attached")
	}
	if seasonGame.Schema.NumMembers != 69 {
		t.Fatalf("SeasonGame schema members = %d, want 69", seasonGame.Schema.NumMembers)
	}

	team, ok := file.GetTableByName("Team")
	if !ok {
		t.Fatal("Team table not found")
	}
	for _, table := range tables {
		if table.Name() != "Team" {
			continue
		}
		if table.AllocatedRecordCount() > team.AllocatedRecordCount() {
			team = &table
		}
	}
	if team.AllocatedRecordCount() != 143 {
		t.Fatalf("Team records = %d, want 143", team.AllocatedRecordCount())
	}
	if team.Schema == nil {
		t.Fatal("Team schema not attached")
	}
}

func TestParseTableHeaderSeasonGame(t *testing.T) {
	savePath := filepath.Join("..", "data", "DYNASTY-TESTSAVE-27")
	if _, err := os.Stat(savePath); err != nil {
		t.Skip("test save not available:", savePath)
	}

	settings := DefaultSettings()
	file, err := Open(savePath, &settings)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Parse(); err != nil {
		t.Fatal(err)
	}

	var seasonGame *Table
	for i := range file.Tables() {
		table := &file.Tables()[i]
		if table.Name() == "SeasonGame" && table.AllocatedRecordCount() == 983 {
			seasonGame = table
			break
		}
	}
	if seasonGame == nil {
		t.Fatal("main SeasonGame table not found")
	}
	if seasonGame.Header.Marker != TableMarkerSPBF {
		t.Fatalf("marker = %q, want SPBF", seasonGame.Header.Marker)
	}
	if seasonGame.Header.Offset != 0x19dfba5 {
		t.Fatalf("offset = %#x, want 0x19dfba5", seasonGame.Header.Offset)
	}
}
