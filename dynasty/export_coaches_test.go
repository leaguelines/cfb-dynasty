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

func TestCoachExportsStaffContractFields(t *testing.T) {
	file := openTestSave(t)
	export, err := file.ExportWithOptions(ExportOptions{
		Sections: ExportSections{Coaches: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(export.Coaches) == 0 {
		t.Fatal("expected coaches")
	}

	withStatus := 0
	withLength := 0
	withYearsRemaining := 0
	for _, coach := range export.Coaches {
		if coach.ContractStatus != "" {
			withStatus++
		}
		if coach.ContractLength != nil {
			withLength++
		}
		if coach.ContractYearsRemaining != nil {
			withYearsRemaining++
		}
	}

	if withStatus == 0 {
		t.Fatal("expected at least one coach with contractStatus")
	}
	if withLength == 0 {
		t.Fatal("expected at least one coach with contractLength")
	}
	if withYearsRemaining == 0 {
		t.Fatal("expected at least one coach with contractYearsRemaining")
	}
}
