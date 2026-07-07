package dynasty

import (
	"os"
	"path/filepath"
	"testing"
)

func openTestSave(t *testing.T) *File {
	t.Helper()
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
	return file
}

func TestExportRecruitsAndStats(t *testing.T) {
	file := openTestSave(t)
	export, err := file.Export()
	if err != nil {
		t.Fatal(err)
	}

	if len(export.Games) == 0 {
		t.Fatal("expected games in export")
	}
	if len(export.Recruits) == 0 {
		t.Fatal("expected recruits in export")
	}
	if len(export.Recruits) > 5000 {
		t.Fatalf("recruits = %d, want active recruit rows only", len(export.Recruits))
	}

	aarons := 0
	for _, recruit := range export.Recruits {
		if recruit.Player != nil && recruit.Player.LastName == "Aarons" {
			aarons++
		}
	}
	if aarons > 10 {
		t.Fatalf("too many recruits linked to Aarons: %d", aarons)
	}
	if len(export.Teams) == 0 {
		t.Fatal("expected teams in export")
	}
	if len(export.Rosters) == 0 {
		t.Fatal("expected rosters in export")
	}
	if len(export.Recruiting) == 0 {
		t.Fatal("expected recruiting in export")
	}
	if len(export.SeasonTeamStats) == 0 {
		t.Fatal("expected season team stats in export")
	}
	if len(export.Coaches) == 0 {
		t.Fatal("expected coaches in export")
	}
	if len(export.LeavingPlayers) == 0 {
		t.Fatal("expected leaving players in export")
	}
	if len(export.Injuries) == 0 {
		t.Fatal("expected injuries in export")
	}
	if len(export.DepthCharts) == 0 {
		t.Fatal("expected depth charts in export")
	}
	if len(export.PlayerAwards) == 0 {
		t.Fatal("expected player awards in export")
	}
	if len(export.StatRecords) == 0 {
		t.Fatal("expected stat records in export")
	}

	withPlayer := 0
	for _, recruit := range export.Recruits {
		if recruit.Player != nil && recruit.Player.FirstName != "" {
			withPlayer++
			if recruit.NationalRank != nil && *recruit.NationalRank > 0 {
				t.Logf("recruit %d %s %s rank=%d stage=%s playerOVR=%v",
					recruit.ID, recruit.Player.FirstName, recruit.Player.LastName,
					*recruit.NationalRank, recruit.RecruitStage, recruit.Player.Overall)
				break
			}
		}
	}
	if withPlayer == 0 {
		t.Fatal("expected at least one recruit with linked player attributes")
	}
}

func TestExportWithOptionsSections(t *testing.T) {
	file := openTestSave(t)

	gamesOnly, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{Games: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(gamesOnly.Games) == 0 {
		t.Fatal("expected games")
	}
	if len(gamesOnly.Recruits) != 0 {
		t.Fatalf("recruits = %d, want 0", len(gamesOnly.Recruits))
	}
	if gamesOnly.Season != nil {
		t.Fatal("expected no season in games-only export")
	}

	recruitsOnly, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{Recruits: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(recruitsOnly.Recruits) == 0 {
		t.Fatal("expected recruits")
	}
	if len(recruitsOnly.Games) != 0 {
		t.Fatalf("games = %d, want 0", len(recruitsOnly.Games))
	}

	teamsOnly, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{Teams: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(teamsOnly.Teams) == 0 {
		t.Fatal("expected teams")
	}
	if len(teamsOnly.Rosters) != 0 {
		t.Fatalf("rosters = %d, want 0", len(teamsOnly.Rosters))
	}

	rostersOnly, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{Rosters: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rostersOnly.Rosters) == 0 {
		t.Fatal("expected rosters")
	}
	if len(rostersOnly.Teams) != 0 {
		t.Fatalf("teams = %d, want 0", len(rostersOnly.Teams))
	}
	if len(rostersOnly.Rosters[0].Players) == 0 {
		t.Fatal("expected players on first roster")
	}
	foundArchetype := false
	foundArchetypeLabel := false
	for _, player := range rostersOnly.Rosters[0].Players {
		if player.Archetype != "" {
			foundArchetype = true
		}
		if player.ArchetypeLabel != "" {
			foundArchetypeLabel = true
		}
	}
	if !foundArchetype {
		t.Fatal("expected at least one roster player with archetype")
	}
	if !foundArchetypeLabel {
		t.Fatal("expected at least one roster player with archetypeLabel")
	}

	recruitingOnly, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{Recruiting: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(recruitingOnly.Recruiting) == 0 {
		t.Fatal("expected recruiting targets")
	}
	if len(recruitingOnly.Recruits) != 0 {
		t.Fatalf("recruits = %d, want 0", len(recruitingOnly.Recruits))
	}
	foundSchool := false
	for _, target := range recruitingOnly.Recruiting {
		if target.TopSchool != nil && target.TopSchool.TeamName != "" {
			foundSchool = true
			break
		}
	}
	if !foundSchool {
		t.Fatal("expected at least one recruiting target with top school")
	}

	seasonStatsOnly, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{SeasonStats: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(seasonStatsOnly.SeasonTeamStats) == 0 {
		t.Fatal("expected team season stats")
	}
	foundWins := false
	for _, teamStats := range seasonStatsOnly.SeasonTeamStats {
		if teamStats.Stats != nil && teamStats.Stats.Wins != nil && *teamStats.Stats.Wins > 0 {
			foundWins = true
			t.Logf("team %s wins=%d", teamStats.TeamName, *teamStats.Stats.Wins)
			break
		}
	}
	if !foundWins {
		t.Fatal("expected at least one team with season wins")
	}

	coachesOnly, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{Coaches: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(coachesOnly.Coaches) == 0 {
		t.Fatal("expected coaches")
	}
	foundHeadCoach := false
	for _, coach := range coachesOnly.Coaches {
		if coach.Position == "HeadCoach" && coach.TeamName != "" {
			foundHeadCoach = true
			break
		}
	}
	if !foundHeadCoach {
		t.Fatal("expected at least one head coach with team")
	}

	leavingOnly, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{LeavingPlayers: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(leavingOnly.LeavingPlayers) == 0 {
		t.Fatal("expected leaving players")
	}
	if leavingOnly.LeavingPlayers[0].FirstName == "" {
		t.Fatal("expected player name on leaving player export")
	}

	injuriesOnly, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{Injuries: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(injuriesOnly.Injuries) == 0 {
		t.Fatal("expected injuries")
	}
	if injuriesOnly.Injuries[0].Type == "" {
		t.Fatal("expected injury type")
	}

	depthOnly, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{DepthCharts: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(depthOnly.DepthCharts) == 0 {
		t.Fatal("expected depth charts")
	}

	historyOnly, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{History: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(historyOnly.PlayerAwards) == 0 {
		t.Fatal("expected player awards")
	}
	if len(historyOnly.StatRecords) == 0 {
		t.Fatal("expected stat records")
	}
}

func TestExportOptionsDefaults(t *testing.T) {
	opts := DefaultExportOptions()
	if !opts.IncludeGames() || !opts.IncludeRecruits() || !opts.IncludeSeason() ||
		!opts.IncludeTeams() || !opts.IncludeRosters() || !opts.IncludeRecruiting() ||
		!opts.IncludeSeasonStats() || !opts.IncludeCoaches() || !opts.IncludeLeavingPlayers() ||
		!opts.IncludeInjuries() || !opts.IncludeDepthCharts() || !opts.IncludeHistory() {
		t.Fatal("default export options should include all sections")
	}
	if !opts.IncludeGameStats() {
		t.Fatal("default export options should include game stats")
	}
	opts.OmitGameStats = true
	if opts.IncludeGameStats() {
		t.Fatal("omit game stats should disable stat export")
	}
}

func TestGameIndexFromStatReference(t *testing.T) {
	ref := &RecordReference{TableID: 0, RowNumber: 101}
	idx, ok := GameIndexFromStatReference(ref, 983)
	if !ok || idx != 100 {
		t.Fatalf("idx=%d ok=%v want 100", idx, ok)
	}
}

func TestRecordByReferenceTeam(t *testing.T) {
	file := openTestSave(t)
	team, ok := file.PrimaryTableByName("Team")
	if !ok {
		t.Fatal("Team not found")
	}
	if err := team.ReadRecords(); err != nil {
		t.Fatal(err)
	}

	ref := &RecordReference{TableID: team.Header.TableID, RowNumber: 80}
	row, ok := file.RecordByReference("Team", ref)
	if !ok {
		t.Fatal("expected team row")
	}
	if stringField(row, "LongName") == "" {
		t.Fatal("expected long name on team row")
	}
}
