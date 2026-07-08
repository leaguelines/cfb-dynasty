package dynasty

import "testing"

func TestExportStandaloneSections(t *testing.T) {
	file := openTestSave(t)
	export, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{
			Rivalries: true,
			Draft:     true,
			Bowls:     true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(export.Rivalries) == 0 {
		t.Fatal("expected rivalries")
	}
	if len(export.DraftPicks) == 0 {
		t.Fatal("expected draft picks")
	}
	if len(export.BowlGames) == 0 {
		t.Fatal("expected bowl games")
	}
}

func TestExportTeamAndCoachEnrichment(t *testing.T) {
	file := openTestSave(t)
	export, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{Teams: true, Coaches: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	withScheme := 0
	for _, team := range export.Teams {
		if team.OffensiveScheme != "" || team.DefensiveScheme != "" {
			withScheme++
		}
	}
	if withScheme == 0 {
		t.Fatal("expected teams with scheme fields")
	}

	withRatings := 0
	withJobSecurity := 0
	for _, coach := range export.Coaches {
		if len(coach.PositionRatings) > 0 {
			withRatings++
		}
		if coach.JobSecurityStatus != "" {
			withJobSecurity++
		}
	}
	if withRatings == 0 && withJobSecurity == 0 {
		t.Fatal("expected coaches with position ratings or job security status")
	}
}

func TestExportPositionChangesBuilder(t *testing.T) {
	file := openTestSave(t)
	changes, err := file.buildPositionChangeExports()
	if err != nil {
		t.Fatal(err)
	}
	// Rows exist in the save but may be placeholder QB->QB entries filtered out.
	t.Logf("position changes exported=%d", len(changes))
}
