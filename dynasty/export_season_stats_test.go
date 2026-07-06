package dynasty

import "testing"

func TestSeasonStatIntOK(t *testing.T) {
	record := Record{Fields: map[string]FieldValue{
		"PASSYARDS": {Int: -16384},
		"PASSTDS":   {Int: 12},
	}}
	if _, ok := seasonStatIntOK(record, "PASSYARDS"); ok {
		t.Fatal("expected negative sentinel to be filtered")
	}
	if v, ok := seasonStatIntOK(record, "PASSTDS"); !ok || v != 12 {
		t.Fatalf("PASSTDS = %v ok=%v", v, ok)
	}
}

func TestBuildTeamSeasonStatsExport(t *testing.T) {
	record := Record{Fields: map[string]FieldValue{
		"WINS":      {Int: 11},
		"LOSSES":    {Int: 2},
		"OFFYARDS":  {Int: 4500},
		"GIVEAWAYS": {Int: 8},
	}}
	stats := buildTeamSeasonStatsExport(record)
	if stats == nil {
		t.Fatal("expected team season stats")
	}
	if stats.Wins == nil || *stats.Wins != 11 {
		t.Fatalf("wins = %v", stats.Wins)
	}
	if stats.TotalYards == nil || *stats.TotalYards != 4500 {
		t.Fatalf("total yards = %v", stats.TotalYards)
	}
}
