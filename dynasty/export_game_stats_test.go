package dynasty

import "testing"

func TestGameStatIntOKIncludesZero(t *testing.T) {
	record := Record{
		Fields: map[string]FieldValue{
			"PASSYARDS": {Int: 0},
			"PASSTDS":   {Int: 2},
			"MISSING":   {Int: -1},
		},
	}
	if v, ok := gameStatIntOK(record, "PASSYARDS"); !ok || v != 0 {
		t.Fatalf("PASSYARDS = %d ok=%v want 0 true", v, ok)
	}
	if v, ok := gameStatIntOK(record, "PASSTDS"); !ok || v != 2 {
		t.Fatalf("PASSTDS = %d ok=%v want 2 true", v, ok)
	}
	if _, ok := gameStatIntOK(record, "MISSING"); ok {
		t.Fatal("expected negative stat to be rejected")
	}
	if _, ok := gameStatIntOK(record, "NOT_PRESENT"); ok {
		t.Fatal("expected missing stat to be rejected")
	}
}

func TestBuildOffensiveGameStatsExportIncludesZeros(t *testing.T) {
	record := Record{
		Fields: map[string]FieldValue{
			"PASSYARDS":      {Int: 0},
			"PASSTDS":        {Int: 0},
			"PASSATTEMPTS":   {Int: 3},
			"PASSCOMPLETED":  {Int: 2},
			"RUSHYARDS":      {Int: 0},
			"RECEIVEYARDS":   {Int: 0},
			"RECEIVECATCHES": {Int: 0},
		},
	}
	stats := buildOffensiveGameStatsExport(record)
	if stats == nil {
		t.Fatal("expected offensive stats export")
	}
	if stats.PassYards == nil || *stats.PassYards != 0 {
		t.Fatalf("passYards = %v want 0", stats.PassYards)
	}
	if stats.PassTDs == nil || *stats.PassTDs != 0 {
		t.Fatalf("passTDs = %v want 0", stats.PassTDs)
	}
	if stats.RushYards == nil || *stats.RushYards != 0 {
		t.Fatalf("rushYards = %v want 0", stats.RushYards)
	}
	if stats.PassAttempts == nil || *stats.PassAttempts != 3 {
		t.Fatalf("passAttempts = %v want 3", stats.PassAttempts)
	}
}

func TestBuildDefensiveGameStatsExportIncludesZeros(t *testing.T) {
	record := Record{
		Fields: map[string]FieldValue{
			"DEFTACKLES":    {Int: 0},
			"ASSDEFTACKLES": {Int: 0},
			"DLINESACKS":    {Int: 1},
		},
	}
	stats := buildDefensiveGameStatsExport(record)
	if stats == nil {
		t.Fatal("expected defensive stats export")
	}
	if stats.Tackles == nil || *stats.Tackles != 0 {
		t.Fatalf("tackles = %v want 0", stats.Tackles)
	}
	if stats.AssistTackles == nil || *stats.AssistTackles != 0 {
		t.Fatalf("assistTackles = %v want 0", stats.AssistTackles)
	}
}
