package dynasty

import (
	"os"
	"path/filepath"
	"testing"
)

func openSeasonSave(t *testing.T) *File {
	t.Helper()
	savePath := filepath.Join("..", "data", "DYNASTY-2026OFFLINEFINAL")
	schemaDir := filepath.Join("..", "data")
	if _, err := os.Stat(savePath); err != nil {
		t.Skip("season save not available:", savePath)
	}

	settings := DefaultSettings()
	settings.SchemaDir = schemaDir
	settings.AutoParse = true
	file, err := Open(savePath, &settings)
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func TestTeamIDFromRecordSkipsFCS(t *testing.T) {
	okRec := Record{Fields: map[string]FieldValue{"TeamIndex": {Int: 125}}}
	id, ok := teamIDFromRecord(okRec)
	if !ok || id != 125 {
		t.Fatalf("TeamIndex 125 = (%d, %v) want (125, true)", id, ok)
	}

	fcs := Record{Fields: map[string]FieldValue{"TeamIndex": {Int: fcsTeamIndexSentinel}}}
	if _, ok := teamIDFromRecord(fcs); ok {
		t.Fatal("expected FCS TeamIndex sentinel to be rejected")
	}
}

func TestTeamMapsRemapRowToCanonicalID(t *testing.T) {
	file := openSeasonSave(t)
	teams := file.teamMaps()

	// Appalachian State sits at a low row with a high TeamIndex — the bug that
	// prompted issue #2. Row 3 must export as 125, not 3.
	id, ok := teams.exportID(3)
	if !ok {
		t.Fatal("expected row 3 (Appalachian State) to remap")
	}
	if id != 125 {
		t.Fatalf("App State exportID = %d want 125", id)
	}
	if name := teams.nameFromID(125); name != "Appalachian State" {
		t.Fatalf("nameFromID(125) = %q want Appalachian State", name)
	}
	if name := teams.nameFromRow(3); name != "Appalachian State" {
		t.Fatalf("nameFromRow(3) = %q want Appalachian State", name)
	}

	// Air Force is TeamIndex 0 at row 0 and must remain addressable.
	id, ok = teams.exportID(0)
	if !ok || id != 0 {
		t.Fatalf("Air Force exportID = (%d, %v) want (0, true)", id, ok)
	}
	if name := teams.nameFromID(0); name != "Air Force" {
		t.Fatalf("nameFromID(0) = %q want Air Force", name)
	}

	// FCS placeholders use TeamIndex 255 and must not appear in the ID map.
	if _, ok := teams.idByRow[30]; ok {
		t.Fatalf("FCS row unexpectedly mapped to id %d", teams.idByRow[30])
	}
}

func TestExportedTeamIDsMatchTeamIndex(t *testing.T) {
	file := openSeasonSave(t)
	exports, err := file.buildTeamExports()
	if err != nil {
		t.Fatal(err)
	}
	byName := make(map[string]int, len(exports))
	for _, team := range exports {
		byName[team.LongName] = team.ID
	}

	// Spot-check against the issue #2 reference table (and Air Force = 0).
	want := map[string]int{
		"Air Force":         0,
		"Akron":             1,
		"Alabama":           2,
		"Arizona":           3,
		"UCF":               17,
		"Appalachian State": 125,
		"Charlotte":         126,
		"Delaware":          134,
		"Missouri State":    135,
		"Sacramento State":  137,
	}
	for name, id := range want {
		got, ok := byName[name]
		if !ok {
			t.Errorf("missing team %q in export", name)
			continue
		}
		if got != id {
			t.Errorf("%s id = %d want %d", name, got, id)
		}
	}
}

func TestRosterTeamIDsMatchTeamIndex(t *testing.T) {
	file := openSeasonSave(t)
	rosters, err := file.buildRosterExports()
	if err != nil {
		t.Fatal(err)
	}
	byName := make(map[string]int, len(rosters))
	for _, roster := range rosters {
		byName[roster.TeamName] = roster.TeamID
		if roster.TeamName == "" {
			t.Fatalf("roster teamId=%d missing name", roster.TeamID)
		}
		for _, player := range roster.Players {
			if player.TeamIndex == nil || *player.TeamIndex != roster.TeamID {
				t.Fatalf("%s player %s %s teamIndex=%v want %d",
					roster.TeamName, player.FirstName, player.LastName, player.TeamIndex, roster.TeamID)
			}
		}
	}
	if byName["Appalachian State"] != 125 {
		t.Fatalf("App State roster teamId = %d want 125", byName["Appalachian State"])
	}
	if byName["Air Force"] != 0 {
		t.Fatalf("Air Force roster teamId = %d want 0", byName["Air Force"])
	}
	if byName["UCF"] != 17 {
		t.Fatalf("UCF roster teamId = %d want 17", byName["UCF"])
	}
	// Spot-check canonical roster IDs for teams involved in row/canonical collisions.
	if byName["Middle Tennessee"] != 51 {
		t.Fatalf("Middle Tennessee roster teamId = %d want 51", byName["Middle Tennessee"])
	}
	if byName["Oklahoma"] != 69 {
		t.Fatalf("Oklahoma roster teamId = %d want 69", byName["Oklahoma"])
	}
	if byName["Temple"] != 88 {
		t.Fatalf("Temple roster teamId = %d want 88", byName["Temple"])
	}

	// Player TeamIndex stores canonical IDs (e.g. Oklahoma = 69).
	okRoster := byName["Oklahoma"]
	for _, roster := range rosters {
		if roster.TeamID != okRoster {
			continue
		}
		for _, player := range roster.Players {
			if player.ID == 58 { // Adebawore, TeamIndex 69
				return
			}
		}
	}
	t.Fatal("expected Oklahoma player id 58 on Oklahoma roster")
}

func TestPlayerTeamIDUsesCanonicalNotRow(t *testing.T) {
	file := openSeasonSave(t)
	teams := file.teamMaps()

	// Row 69 is Middle Tennessee (canonical 51); player TeamIndex 69 is Oklahoma.
	if id, ok := teams.exportID(69); !ok || id != 51 {
		t.Fatalf("exportID(69) = (%d, %v) want (51, true)", id, ok)
	}
	if id, ok := teams.playerTeamID(69); !ok || id != 69 {
		t.Fatalf("playerTeamID(69) = (%d, %v) want (69, true)", id, ok)
	}
	if name := teams.nameFromID(69); name != "Oklahoma" {
		t.Fatalf("nameFromID(69) = %q want Oklahoma", name)
	}
}

func TestJohnMateerOnOklahomaRoster(t *testing.T) {
	file := openSeasonSave(t)
	rosters, err := file.buildRosterExports()
	if err != nil {
		t.Fatal(err)
	}
	var oklahoma, temple *RosterExport
	for i := range rosters {
		switch rosters[i].TeamName {
		case "Oklahoma":
			oklahoma = &rosters[i]
		case "Temple":
			temple = &rosters[i]
		}
	}
	if oklahoma == nil {
		t.Fatal("missing Oklahoma roster")
	}
	if oklahoma.TeamID != 69 {
		t.Fatalf("Oklahoma teamId = %d want 69", oklahoma.TeamID)
	}

	var mateerOnOklahoma bool
	for _, player := range oklahoma.Players {
		if player.FirstName == "John" && player.LastName == "Mateer" {
			mateerOnOklahoma = true
			if player.Position != "QB" {
				t.Fatalf("Mateer position = %q want QB", player.Position)
			}
			if player.TeamIndex == nil || *player.TeamIndex != 69 {
				t.Fatalf("Mateer teamIndex = %v want 69", player.TeamIndex)
			}
		}
	}
	if !mateerOnOklahoma {
		t.Fatal("John Mateer not found on Oklahoma roster")
	}

	if temple != nil {
		for _, player := range temple.Players {
			if player.FirstName == "John" && player.LastName == "Mateer" {
				t.Fatal("John Mateer incorrectly listed on Temple roster")
			}
			if player.FirstName == "Scotty" && player.LastName == "Fox Jr." {
				if player.TeamIndex == nil || *player.TeamIndex != 88 {
					t.Fatalf("Temple QB teamIndex = %v want 88", player.TeamIndex)
				}
			}
		}
	}

	// Temple canonical id 88 must not be remapped onto Oklahoma via exportID.
	for _, player := range oklahoma.Players {
		if player.FirstName == "Scotty" && player.LastName == "Fox Jr." {
			t.Fatal("Temple QB Scotty Fox Jr. incorrectly on Oklahoma roster")
		}
	}
}
