package dynasty

import "testing"

func TestBuildCareerStatsIndexSkipsNilCareerReference(t *testing.T) {
	skipIfShortIntegration(t)
	file := openTestSave(t)
	idx, err := file.buildCareerStatsIndex()
	if err != nil {
		t.Fatal(err)
	}

	playerTable, _ := file.PrimaryTableByName("Player")
	playerTable.ReadRecords()
	careerOff, _ := file.PrimaryTableByName("CareerOffensiveStats")
	careerOff.ReadRecords()

	nullRefPlayers := 0
	for _, player := range playerTable.Records {
		if player.Index >= int(playerTable.ActiveRecordCount()) {
			break
		}
		ref, ok := player.Get("CareerStats")
		if !ok || !isNilReference(ref.Reference) {
			continue
		}
		nullRefPlayers++
		if row, mapped := idx.offByPlayer[player.Index]; mapped {
			t.Fatalf("player %d with nil CareerStats mapped to career row %d", player.Index, row)
		}
	}
	if nullRefPlayers == 0 {
		t.Fatal("expected players with nil CareerStats in test save")
	}
}

func TestPlayerCareerStatsFromRecordsIncludesCareerProduction(t *testing.T) {
	stats := playerCareerStatsFromRecords(Record{Fields: map[string]FieldValue{
		"GAMESPLAYED":  {Int: 144},
		"GAMESSTARTED": {Int: 140},
		"DOWNSPLAYED":  {Int: 4200},
		"PASSYARDS":    {Int: 33946},
		"PASSTDS":      {Int: 267},
	}}, Record{}, nil)
	if stats == nil {
		t.Fatal("expected career stats")
	}
	if stats.GamesPlayed == nil || *stats.GamesPlayed != 144 {
		t.Fatalf("gamesPlayed = %v", stats.GamesPlayed)
	}
	if stats.DownsPlayed == nil || *stats.DownsPlayed != 4200 {
		t.Fatalf("downsPlayed = %v", stats.DownsPlayed)
	}
	if stats.Offense == nil || stats.Offense.PassYards == nil || *stats.Offense.PassYards != 33946 {
		t.Fatalf("offense = %#v", stats.Offense)
	}
}

func TestSelectCareerStatRecordsPrefersSummedSeasonRows(t *testing.T) {
	careerOff := Record{Fields: map[string]FieldValue{
		"GAMESPLAYED": {Int: 10},
		"PASSYARDS":   {Int: 2191},
	}}
	seasonRows := []Record{
		{Fields: map[string]FieldValue{
			"GAMESPLAYED": {Int: 12},
			"PASSYARDS":   {Int: 3100},
		}},
		{Fields: map[string]FieldValue{
			"GAMESPLAYED": {Int: 11},
			"PASSYARDS":   {Int: 2800},
		}},
		{Fields: map[string]FieldValue{
			"GAMESPLAYED": {Int: 13},
			"PASSYARDS":   {Int: 3500},
		}},
	}
	table := &Table{Records: seasonRows}

	offRecord, _ := selectCareerStatRecords(careerOff, Record{}, careerStatsIndex{
		seasonOffByPlayer: map[int][]int{42: {0, 1, 2}},
		seasonOffTable:    table,
	}, 42)

	if gp, ok := intFieldOK(offRecord, "GAMESPLAYED"); !ok || gp != 36 {
		t.Fatalf("merged gamesPlayed = %v ok=%v", gp, ok)
	}
	if py, ok := intFieldOK(offRecord, "PASSYARDS"); !ok || py != 9400 {
		t.Fatalf("merged passYards = %v ok=%v", py, ok)
	}
}

func TestSeasonStatRowKeysGroupsBySlot(t *testing.T) {
	keys := seasonStatRowKeys(
		[]int{10, 25},
		[]int{11, 25},
		map[int]int{10: 1, 25: 3},
		map[int]int{11: 2, 25: 3},
	)
	if len(keys) != 3 {
		t.Fatalf("keys = %#v", keys)
	}
	if keys[0].slot != 1 || keys[0].offRow != 10 || keys[0].defRow != 0 {
		t.Fatalf("slot1 = %#v", keys[0])
	}
	if keys[2].slot != 3 || keys[2].offRow != 25 || keys[2].defRow != 25 {
		t.Fatalf("slot3 = %#v", keys[2])
	}
}
