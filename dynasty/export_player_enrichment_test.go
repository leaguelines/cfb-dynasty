package dynasty

import "testing"

func TestPhysicalAbilitiesFromRecord(t *testing.T) {
	abilities := physicalAbilitiesFromRecord(Record{Fields: map[string]FieldValue{
		"PhysicalAbility1": {String: "Gold"},
		"PhysicalAbility2": {String: "None"},
		"PhysicalAbility3": {String: "Silver_"},
	}})
	if len(abilities) != 2 {
		t.Fatalf("abilities = %#v, want 2", abilities)
	}
	if abilities[0].Slot != 1 || abilities[0].Tier != "Gold" {
		t.Fatalf("slot1 = %#v", abilities[0])
	}
	if abilities[1].Tier != "Silver" {
		t.Fatalf("slot3 = %#v", abilities[1])
	}
}

func TestMentalAbilitiesFromRecord(t *testing.T) {
	rank := 2
	abilities := mentalAbilitiesFromRecord(Record{Fields: map[string]FieldValue{
		"MentalAbility1":     {String: "WinningTime"},
		"MentalAbilityRank1": {Int: 2},
		"MentalAbility2":     {String: "None"},
	}})
	if len(abilities) != 1 {
		t.Fatalf("abilities = %#v", abilities)
	}
	if abilities[0].Name != "WinningTime" || abilities[0].Rank == nil || *abilities[0].Rank != rank {
		t.Fatalf("ability = %#v", abilities[0])
	}
}

func TestPlayerCareerStatsFromRecord(t *testing.T) {
	stats := playerCareerStatsFromRecords(Record{Fields: map[string]FieldValue{
		"GAMESPLAYED":  {Int: 14},
		"GAMESSTARTED": {Int: 12},
		"DOWNSPLAYED":  {Int: 80},
		"PASSYARDS":    {Int: 3200},
	}}, Record{}, nil)
	if stats == nil {
		t.Fatal("expected career stats")
	}
	if stats.GamesPlayed == nil || *stats.GamesPlayed != 14 {
		t.Fatalf("gamesPlayed = %v", stats.GamesPlayed)
	}
	if stats.Offense == nil || stats.Offense.PassYards == nil || *stats.Offense.PassYards != 3200 {
		t.Fatalf("offense = %#v", stats.Offense)
	}
}

func TestArchetypeTraitsFromRecord(t *testing.T) {
	traits := archetypeTraitsFromRecord(Record{Fields: map[string]FieldValue{
		"PT_QBPOCKETPASSER": {Bool: true},
		"PT_AGGRESSIVERECEIVER": {Bool: true},
	}})
	if len(traits) != 2 {
		t.Fatalf("traits = %#v", traits)
	}
	foundPocket := false
	foundAggressive := false
	for _, trait := range traits {
		switch trait {
		case "Pocket Passer":
			foundPocket = true
		case "AGGRESSIVERECEIVER":
			foundAggressive = true
		}
	}
	if !foundPocket || !foundAggressive {
		t.Fatalf("traits = %#v", traits)
	}
}

func TestApplyPlayerEnrichment(t *testing.T) {
	skipIfShortIntegration(t)
	file := openTestSave(t)
	teams := file.teamMaps()
	record := Record{Fields: map[string]FieldValue{
		"FirstName":              {String: "Test"},
		"LastName":               {String: "Player"},
		"RedshirtStatus":         {String: "Eligible_"},
		"IsNIL":                  {Bool: true},
		"BaseNILValue":           {Int: 100},
		"CurrentNILCompensation": {Int: 150},
		"Motivation1":            {String: "PlayingTime_"},
		"PhysicalAbility1":       {String: "Gold"},
		"MentalAbility1":         {String: "TheNatural"},
		"MentalAbilityRank1":     {Int: 1},
		"WasPreviouslyInjured":   {Bool: true},
	}}
	player := buildPlayerExport(record)
	applyPlayerEnrichment(file, player, record, teams, careerStatsIndex{})
	if player.RedshirtStatus != "Eligible" {
		t.Fatalf("redshirtStatus = %q", player.RedshirtStatus)
	}
	if !player.IsNIL || player.NILBaseValue == nil || *player.NILBaseValue != 100 {
		t.Fatalf("nil = %#v", player)
	}
	if len(player.Motivations) != 1 || player.Motivations[0] != "PlayingTime" {
		t.Fatalf("motivations = %#v", player.Motivations)
	}
	if len(player.PhysicalAbilities) != 1 {
		t.Fatalf("physical = %#v", player.PhysicalAbilities)
	}
	if !player.WasPreviouslyInjured {
		t.Fatal("expected wasPreviouslyInjured")
	}
}

func TestExportRosters_PlayerEnrichment(t *testing.T) {
	skipIfShortIntegration(t)
	file := openTestSave(t)
	export, err := file.ExportWithOptions(ExportOptions{Sections: ExportSections{Rosters: true}})
	if err != nil {
		t.Fatal(err)
	}

	withNIL, withMotivation, withPhysical, withCareer := 0, 0, 0, 0
	for _, roster := range export.Rosters {
		for _, player := range roster.Players {
			if player.IsNIL {
				withNIL++
			}
			if len(player.Motivations) > 0 {
				withMotivation++
			}
			if len(player.PhysicalAbilities) > 0 {
				withPhysical++
			}
			if player.CareerStats != nil {
				withCareer++
			}
		}
	}
	if withMotivation == 0 {
		t.Fatal("expected players with motivations")
	}
	if withPhysical == 0 {
		t.Fatal("expected players with physical abilities")
	}
	t.Logf("NIL=%d motivations=%d physical=%d career=%d", withNIL, withMotivation, withPhysical, withCareer)
}

func TestExportRecruits_PlayerEnrichment(t *testing.T) {
	skipIfShortIntegration(t)
	file := openTestSave(t)
	export, err := file.ExportWithOptions(ExportOptions{Sections: ExportSections{Recruits: true}})
	if err != nil {
		t.Fatal(err)
	}
	withDealbreaker := 0
	for _, recruit := range export.Recruits {
		if recruit.Player == nil {
			continue
		}
		if recruit.Player.RecruitingDealbreaker != "" {
			withDealbreaker++
		}
	}
	if withDealbreaker == 0 {
		t.Fatal("expected recruits with recruiting dealbreaker")
	}
}
