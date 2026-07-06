package dynasty

import "testing"

func TestCoachCareerStatsFromRecord(t *testing.T) {
	record := Record{Fields: map[string]FieldValue{
		"Wins":      {Int: 42},
		"Losses":    {Int: 18},
		"BowlWins":  {Int: 3},
		"NCWins":    {Int: 1},
	}}
	stats := coachCareerStatsFromRecord(record)
	if stats == nil {
		t.Fatal("expected career stats")
	}
	if stats.Wins == nil || *stats.Wins != 42 {
		t.Fatalf("wins = %v", stats.Wins)
	}
	if stats.BowlWins == nil || *stats.BowlWins != 3 {
		t.Fatalf("bowl wins = %v", stats.BowlWins)
	}
}
