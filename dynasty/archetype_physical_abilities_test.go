package dynasty

import "testing"

func TestPhysicalAbilityName(t *testing.T) {
	tests := []struct {
		position string
		archetype string
		slot int
		want string
	}{
		{"QB", "Pocket Passer", 1, "Resistance"},
		{"QB", "Pocket Passer", 4, "Dot!"},
		{"WR", "Gritty Possession", 1, "Second Level"},
		{"TE", "Gritty Possession", 1, "Workhorse"},
		{"LE", "Edge Setter", 3, "Outside Disruptor"},
		{"DT", "Gap Specialist", 3, "Outside Disruptor"},
		{"K", "Accurate", 1, "Chip Shot"},
		{"P", "Power", 3, "Coffin Corner"},
		{"OT", "Well-Rounded", 2, "Outside Shield"},
		{"QB", "Unknown Archetype", 1, ""},
	}
	for _, tt := range tests {
		if got := physicalAbilityName(tt.position, tt.archetype, tt.slot); got != tt.want {
			t.Fatalf("physicalAbilityName(%q, %q, %d) = %q, want %q", tt.position, tt.archetype, tt.slot, got, tt.want)
		}
	}
}

func TestPhysicalAbilitiesFromRecordIncludesNames(t *testing.T) {
	abilities := physicalAbilitiesFromRecord(Record{Fields: map[string]FieldValue{
		"Position":         {String: "QB"},
		"PT_QBPOCKETPASSER": {Bool: true},
		"PhysicalAbility1": {String: "Gold"},
		"PhysicalAbility2": {String: "Silver_"},
	}})
	if len(abilities) != 2 {
		t.Fatalf("abilities = %#v, want 2", abilities)
	}
	if abilities[0].Name != "Resistance" || abilities[1].Name != "Step Up" {
		t.Fatalf("ability names = %#v", abilities)
	}
}
