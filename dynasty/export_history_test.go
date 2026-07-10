package dynasty

import "testing"

func TestHistoryExportsSkipInactiveRows(t *testing.T) {
	skipIfShortIntegration(t)
	file := openSeasonSave(t)
	export, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{History: true},
	})
	if err != nil {
		t.Fatal(err)
	}

	awardTable, _ := file.PrimaryTableByName("LeagueHistoryAward")
	if len(export.LeagueAwards) != int(awardTable.ActiveRecordCount()) {
		t.Fatalf("leagueAwards = %d active = %d", len(export.LeagueAwards), awardTable.ActiveRecordCount())
	}

	playerAwardTable, _ := file.PrimaryTableByName("PlayerAward")
	if len(export.PlayerAwards) > int(playerAwardTable.ActiveRecordCount()) {
		t.Fatalf("playerAwards = %d exceeds active %d", len(export.PlayerAwards), playerAwardTable.ActiveRecordCount())
	}

	champTable, _ := file.PrimaryTableByName("LeagueHistoryConferenceChampion")
	if len(export.ConferenceChampions) != int(champTable.ActiveRecordCount()) {
		t.Fatalf("conferenceChampions = %d active = %d", len(export.ConferenceChampions), champTable.ActiveRecordCount())
	}

	if len(export.RecordBook) == 0 {
		t.Fatal("expected a populated record book")
	}

	for _, award := range export.LeagueAwards {
		if award.FirstName == "Curt" && award.LastName == "Curt" && award.TeamDisplayName == "Curt" {
			t.Fatalf("placeholder league award at id %d: %+v", award.ID, award)
		}
	}
}

func TestStatRecordPositionsInferredFromStatType(t *testing.T) {
	skipIfShortIntegration(t)
	file := openSeasonSave(t)
	export, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{History: true},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Record rows store a real position for the holder, but league ranks 2..N
	// leave it at the schema default (QB). Inference must fix categories that a
	// quarterback can never lead; positions like TE catches, LB sacks, or a
	// mobile QB's rushing record are legitimately kept as stored.
	neverQB := map[string]bool{
		"ReceiveCatches": true, "ReceiveYards": true, "ReceiveTDs": true,
		"DefensiveInts": true, "DefensiveSacks": true,
	}
	for _, rec := range export.RecordBook {
		if rec.Position == "QB" && neverQB[rec.StatType] {
			t.Errorf("%s scope %s %s: position=QB on non-QB category (default not inferred)",
				rec.LastName, rec.Scope, rec.StatType)
		}
	}
}

func TestRecordBookScopesAndPeriods(t *testing.T) {
	skipIfShortIntegration(t)
	file := openSeasonSave(t)
	export, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{History: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(export.RecordBook) == 0 {
		t.Fatal("expected record book entries")
	}

	scopes := map[string]int{}
	periods := map[string]int{}
	teamNames := map[string]struct{}{}
	confNames := map[string]struct{}{}
	leagueTopRank := 0
	for _, e := range export.RecordBook {
		scopes[e.Scope]++
		periods[e.Period]++
		switch e.Scope {
		case "team":
			if e.ScopeName == "" {
				t.Errorf("team record entry missing scopeName: %+v", e)
			}
			teamNames[e.ScopeName] = struct{}{}
		case "conference":
			confNames[e.ScopeName] = struct{}{}
		case "league":
			if e.Rank > leagueTopRank {
				leagueTopRank = e.Rank
			}
		}
		if e.Rank < 1 {
			t.Errorf("entry with invalid rank: %+v", e)
		}
	}

	for _, scope := range []string{"league", "conference", "team"} {
		if scopes[scope] == 0 {
			t.Errorf("no %s record entries", scope)
		}
	}
	for _, period := range []string{"career", "season", "game"} {
		if periods[period] == 0 {
			t.Errorf("no %s record entries", period)
		}
	}
	// The record book should cover many teams and multiple conferences.
	if len(teamNames) < 50 {
		t.Errorf("expected records for many teams, got %d", len(teamNames))
	}
	if len(confNames) < 2 {
		t.Errorf("expected records for multiple conferences, got %d", len(confNames))
	}
	// League boards are ranked top-N, so rank should exceed 1 somewhere.
	if leagueTopRank < 2 {
		t.Errorf("expected league boards to be ranked (top rank %d)", leagueTopRank)
	}
}

func TestPositionGroupFromStatType(t *testing.T) {
	tests := []struct {
		statType string
		want     string
		ok       bool
	}{
		{"ReceiveCatches", "WR", true},
		{"DefensiveInts", "DB", true},
		{"DefensiveSacks", "DL", true},
		{"PassTds", "QB", true},
		{"Pancakes", "OL", true},
		{"GamesPlayed", "", false},
	}
	for _, tc := range tests {
		got, ok := positionGroupFromStatType(tc.statType)
		if got != tc.want || ok != tc.ok {
			t.Errorf("positionGroupFromStatType(%q) = %q, %v want %q, %v", tc.statType, got, ok, tc.want, tc.ok)
		}
	}
}
