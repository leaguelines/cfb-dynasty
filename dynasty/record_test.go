package dynasty

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTeamShortNameOffset(t *testing.T) {
	savePath := filepath.Join("..", "data", "DYNASTY-TESTSAVE-27")
	schemaDir := filepath.Join("..", "data")
	if _, err := os.Stat(savePath); err != nil {
		t.Skip("test save not available:", savePath)
	}

	settings := DefaultSettings()
	settings.SchemaDir = schemaDir
	settings.AutoParse = true
	file, err := Open(savePath, &settings)
	if err != nil {
		t.Fatal(err)
	}

	team, ok := file.PrimaryTableByName("Team")
	if !ok {
		t.Fatal("Team not found")
	}

	offsets, err := buildOffsetTable(team.Data, team.Schema, team.Header)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range offsets {
		if entry.Name == "ShortName" {
			t.Logf("ShortName: indexOffset=%d offset=%d length=%d secondTable=%v",
				entry.IndexOffset, entry.Offset, entry.Length, entry.ValueInSecondTable)
		}
	}

	if err := team.ReadRecords(); err != nil {
		t.Fatal(err)
	}
	short := stringField(team.Records[0], "ShortName")
	if short == "" {
		t.Fatal("expected ShortName on team 0")
	}
	long := stringField(team.Records[0], "LongName")
	if long == "" {
		t.Fatal("expected LongName on team 0")
	}
	t.Logf("team0 ShortName=%q LongName=%q DisplayName=%q", short, long, stringField(team.Records[0], "DisplayName"))

	if len(team.Records) <= 80 {
		t.Fatal("expected at least 81 team rows")
	}
	if got := stringField(team.Records[80], "LongName"); got != "North Carolina" {
		t.Fatalf("team 80 LongName = %q, want North Carolina", got)
	}

	seasonGame, ok := file.PrimaryTableByName("SeasonGame")
	if !ok {
		t.Fatal("SeasonGame not found")
	}
	if err := seasonGame.ReadRecords(); err != nil {
		t.Fatal(err)
	}
	if len(seasonGame.Records) == 0 {
		t.Fatal("expected SeasonGame records")
	}
	status := stringField(seasonGame.Records[0], "GameStatus")
	if status == "" || status == "First_" {
		t.Fatalf("unexpected GameStatus %q", status)
	}
	t.Logf("game0 GameStatus=%q", status)
}
