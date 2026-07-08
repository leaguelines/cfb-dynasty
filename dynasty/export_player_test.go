package dynasty

import "testing"

func TestPlayerIsAth(t *testing.T) {
	tests := []struct {
		alt1, alt2 string
		want       bool
	}{
		{alt1: "SS", alt2: "Invalid_", want: true},
		{alt1: "Invalid_", alt2: "WR", want: true},
		{alt1: "CB", alt2: "FS", want: true},
		{alt1: "Invalid_", alt2: "Invalid_", want: false},
		{alt1: "", alt2: "", want: false},
		{alt1: "QB", alt2: "", want: true},
	}
	for _, tt := range tests {
		if got := playerIsAth(tt.alt1, tt.alt2); got != tt.want {
			t.Errorf("playerIsAth(%q, %q) = %v, want %v", tt.alt1, tt.alt2, got, tt.want)
		}
	}
}

func TestBuildPlayerExport_IsAthFromAlternates(t *testing.T) {
	player := buildPlayerExport(Record{Fields: map[string]FieldValue{
		"FirstName": {String: "Test"},
		"LastName":  {String: "Athlete"},
		"Position":  {String: "QB"},
	}})
	if player == nil {
		t.Fatal("expected player export")
	}
	applyPlayerAth(player, "SS", "Invalid_")
	if !player.IsAth {
		t.Fatal("expected isAth=true when alternate position is set")
	}

	applyPlayerAth(player, "Invalid_", "Invalid_")
	if player.IsAth {
		t.Fatal("expected isAth=false when no alternate positions are set")
	}
}

func TestExportRecruits_IsAth(t *testing.T) {
	file := openTestSave(t)
	export, err := file.Export()
	if err != nil {
		t.Fatal(err)
	}

	athRecruits := 0
	for _, recruit := range export.Recruits {
		if recruit.Player == nil {
			continue
		}
		hasAlt := isSetAlternatePosition(recruit.AlternatePosition1) ||
			isSetAlternatePosition(recruit.AlternatePosition2)
		if hasAlt != recruit.Player.IsAth {
			t.Fatalf("recruit %d: isAth=%v, want %v (alt1=%q alt2=%q)",
				recruit.ID, recruit.Player.IsAth, hasAlt,
				recruit.AlternatePosition1, recruit.AlternatePosition2)
		}
		if recruit.Player.IsAth {
			athRecruits++
		}
	}
	if athRecruits == 0 {
		t.Fatal("expected at least one ATH recruit in test save")
	}
	t.Logf("ATH recruits=%d", athRecruits)
}

func TestPlayerAlternatePositions(t *testing.T) {
	file := openTestSave(t)
	positions, err := file.playerAlternatePositions()
	if err != nil {
		t.Fatal(err)
	}
	if len(positions) == 0 {
		t.Fatal("expected recruit-to-player alternate position map")
	}

	export, err := file.Export()
	if err != nil {
		t.Fatal(err)
	}
	for _, recruit := range export.Recruits {
		if recruit.Player == nil || !recruit.Player.IsAth {
			continue
		}
		alt, ok := positions[recruit.Player.ID]
		if !ok {
			t.Fatalf("player %d missing from alternate position map", recruit.Player.ID)
		}
		if !playerIsAth(alt[0], alt[1]) {
			t.Fatalf("player %d alt=%v, want ATH", recruit.Player.ID, alt)
		}
		return
	}
	t.Fatal("expected at least one ATH recruit")
}
