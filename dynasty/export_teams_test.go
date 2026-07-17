package dynasty

import "testing"

func TestTeamExportsIncludeOFFDEFOVRRatings(t *testing.T) {
	skipIfShortIntegration(t)
	file := openSeasonSave(t)
	exports, err := file.buildTeamExports()
	if err != nil {
		t.Fatal(err)
	}

	byName := make(map[string]TeamExport, len(exports))
	for _, team := range exports {
		byName[team.LongName] = team
	}

	georgia, ok := byName["Georgia"]
	if !ok {
		t.Fatal("Georgia missing from team export")
	}
	if georgia.OffensiveRating == nil || *georgia.OffensiveRating != 94 {
		t.Fatalf("Georgia offensiveRating = %v want 94", georgia.OffensiveRating)
	}
	if georgia.DefensiveRating == nil || *georgia.DefensiveRating != 96 {
		t.Fatalf("Georgia defensiveRating = %v want 96", georgia.DefensiveRating)
	}
	if georgia.OverallRating == nil || *georgia.OverallRating != 94 {
		t.Fatalf("Georgia overallRating = %v want 94", georgia.OverallRating)
	}

	rated := 0
	for _, team := range exports {
		if team.OverallRating != nil && *team.OverallRating > 0 {
			rated++
		}
	}
	if rated < 100 {
		t.Fatalf("teams with overallRating = %d, want at least 100", rated)
	}
}
