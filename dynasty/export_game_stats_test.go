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

func TestBuildSpecialTeamsStatsExport(t *testing.T) {
	// Row with kick-return activity: all fields populated, including zeros.
	record := Record{
		Fields: map[string]FieldValue{
			"KRETATTEMPTS": {Int: 3},
			"KRETYARDS":    {Int: 142},
			"KRETLONGEST":  {Int: 99},
			"KRETTDS":      {Int: 1},
			"PRETATTEMPTS": {Int: 0},
			"PRETYARDS":    {Int: 0},
			"PRETLONGEST":  {Int: 0},
			"PRETTDS":      {Int: 0},
		},
	}
	stats := buildSpecialTeamsStatsExport(record)
	if stats == nil {
		t.Fatal("expected special teams stats export")
	}
	if stats.KickReturnTDs == nil || *stats.KickReturnTDs != 1 {
		t.Fatalf("kickReturnTDs = %v want 1", stats.KickReturnTDs)
	}
	if stats.KickReturnYards == nil || *stats.KickReturnYards != 142 {
		t.Fatalf("kickReturnYards = %v want 142", stats.KickReturnYards)
	}
	if stats.PuntReturns == nil || *stats.PuntReturns != 0 {
		t.Fatalf("puntReturns = %v want 0", stats.PuntReturns)
	}

	// Row with no return activity at all yields nil (non-returners stay clean).
	empty := Record{
		Fields: map[string]FieldValue{
			"KRETATTEMPTS": {Int: 0},
			"KRETYARDS":    {Int: 0},
			"KRETTDS":      {Int: 0},
			"PRETATTEMPTS": {Int: 0},
			"PRETYARDS":    {Int: 0},
			"PRETTDS":      {Int: 0},
		},
	}
	if got := buildSpecialTeamsStatsExport(empty); got != nil {
		t.Fatalf("expected nil for inactive returner, got %+v", got)
	}
}

func TestMergeSpecialTeams(t *testing.T) {
	kr, pr := 5, 3
	a := &SpecialTeamsStatsExport{KickReturns: &kr}
	b := &SpecialTeamsStatsExport{PuntReturns: &pr}
	merged := mergeSpecialTeams(a, b)
	if merged.KickReturns == nil || *merged.KickReturns != 5 {
		t.Fatalf("kickReturns = %v want 5", merged.KickReturns)
	}
	if merged.PuntReturns == nil || *merged.PuntReturns != 3 {
		t.Fatalf("puntReturns = %v want 3", merged.PuntReturns)
	}
	if got := mergeSpecialTeams(nil, b); got != b {
		t.Fatal("merge(nil, b) should return b")
	}
	if got := mergeSpecialTeams(a, nil); got != a {
		t.Fatal("merge(a, nil) should return a")
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
