package dynasty

import (
	"slices"
	"testing"
)

func TestSkillGroupSlotsFromRecord(t *testing.T) {
	record := Record{Fields: map[string]FieldValue{
		"SkillGroupCap1": {Int: 17},
		"SkillGroupCap2": {Int: 19},
		"SkillGroupCap3": {Int: 15},
		"SkillGroupCap4": {Int: 18},
		"SkillGroupCap5": {Int: 14},
		"SkillGroupCap6": {Int: 16},
	}}
	slots, ok := skillGroupSlotsFromRecord(record)
	if !ok {
		t.Fatal("expected ok")
	}
	wantUnlocked := []int{17, 19, 15, 18, 14, 16}
	wantCapped := []int{3, 1, 5, 2, 6, 4}
	for i := range wantUnlocked {
		if slots.unlocked[i] != wantUnlocked[i] {
			t.Fatalf("unlocked[%d] = %d, want %d", i, slots.unlocked[i], wantUnlocked[i])
		}
		if slots.capped[i] != wantCapped[i] {
			t.Fatalf("capped[%d] = %d, want %d", i, slots.capped[i], wantCapped[i])
		}
	}
	if totalCappedSlots(slots) != 21 {
		t.Fatalf("total capped = %d, want 21", totalCappedSlots(slots))
	}
	if sumInts(slots.unlocked) != 99 {
		t.Fatalf("total unlocked = %d, want 99", sumInts(slots.unlocked))
	}
}

func TestSkillGroupSlotsFromRecord_MissingOrEmpty(t *testing.T) {
	if _, ok := skillGroupSlotsFromRecord(Record{Fields: map[string]FieldValue{"SkillGroupCap1": {Int: 5}}}); ok {
		t.Fatal("expected not ok when caps are incomplete")
	}
	zero := Record{Fields: map[string]FieldValue{}}
	for i := 1; i <= skillGroupCapCount; i++ {
		zero.Fields[skillGroupCapField(i)] = FieldValue{Int: 0}
	}
	if _, ok := skillGroupSlotsFromRecord(zero); ok {
		t.Fatal("expected not ok when all caps are zero")
	}
}

func TestSkillGroupOverall(t *testing.T) {
	rating := func(value int) *int { return &value }
	tests := []struct {
		name       string
		attributes []SkillGroupAttributeExport
		want       int
		wantOK     bool
	}{
		{
			name: "weighted tiers",
			attributes: []SkillGroupAttributeExport{
				{Rating: rating(80), Tier: "Primary"},
				{Rating: rating(100), Tier: "Secondary"},
				{Rating: rating(70), Tier: "Tertiary"},
			},
			want: 82, wantOK: true,
		},
		{
			name: "tiers absent",
			attributes: []SkillGroupAttributeExport{
				{Rating: rating(80)},
				{Rating: rating(100)},
			},
			want: 90, wantOK: true,
		},
		{
			name: "tiers incomplete",
			attributes: []SkillGroupAttributeExport{
				{Rating: rating(80), Tier: "Primary"},
				{Rating: rating(100)},
			},
			want: 90, wantOK: true,
		},
		{
			name:       "rating missing",
			attributes: []SkillGroupAttributeExport{{Tier: "Primary"}},
			wantOK:     false,
		},
		{name: "attributes missing", wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := SkillGroupOverall(tt.attributes)
			if got != tt.want || ok != tt.wantOK {
				t.Fatalf("SkillGroupOverall() = (%d, %v), want (%d, %v)", got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestSkillGroupCurrentLevel(t *testing.T) {
	families := []struct {
		name      string
		positions []string
		minimums  []int
	}{
		{
			name:      "FB and TE",
			positions: []string{"FB", "TE"},
			minimums:  []int{59, 61, 63, 65, 67, 69, 71, 73, 75, 77, 79, 81, 83, 84, 85, 86, 87, 88, 89},
		},
		{
			name:      "WR MIKE K and P",
			positions: []string{"WR", "MIKE", "MLB", "K", "P"},
			minimums:  []int{61, 64, 66, 69, 71, 74, 76, 79, 81, 84, 86, 89, 91, 93, 95, 96, 97, 98, 99},
		},
		{
			name:      "default",
			positions: []string{"QB", "ATH", "te", ""},
			minimums:  []int{63, 65, 67, 69, 71, 73, 75, 77, 79, 81, 83, 85, 87, 89, 91, 93, 95, 97, 99},
		},
	}

	for _, family := range families {
		t.Run(family.name, func(t *testing.T) {
			for _, position := range family.positions {
				for i, minimum := range family.minimums {
					if got, want := skillGroupCurrentLevel(position, minimum-1, 20), i+1; got != want {
						t.Errorf("skillGroupCurrentLevel(%q, %d, 20) = %d, want %d", position, minimum-1, got, want)
					}
					if got, want := skillGroupCurrentLevel(position, minimum, 20), i+2; got != want {
						t.Errorf("skillGroupCurrentLevel(%q, %d, 20) = %d, want %d", position, minimum, got, want)
					}
				}
			}
		})
	}

	if got := skillGroupCurrentLevel("FB", 89, 17); got != 17 {
		t.Errorf("clamped level = %d, want 17", got)
	}
	for _, ceiling := range []int{0, -1} {
		if got := skillGroupCurrentLevel("QB", 99, ceiling); got != 0 {
			t.Errorf("level with ceiling %d = %d, want 0", ceiling, got)
		}
	}
}

func TestSkillGroupCurrentLevelsRequiresSixCompleteGroups(t *testing.T) {
	rating := 99
	groups := make([]SkillGroupExport, skillGroupCapCount)
	for i := range groups {
		groups[i].Attributes = []SkillGroupAttributeExport{{Rating: &rating, Tier: "Primary"}}
	}
	ceilings := []int{20, 19, 18, 17, 16, 0}
	if got := skillGroupCurrentLevels("QB", groups, ceilings); !slices.Equal(got, ceilings) {
		t.Fatalf("levels = %v, want %v", got, ceilings)
	}

	rating = 61
	for _, tt := range []struct {
		position string
		want     []int
	}{
		{position: "FB", want: []int{3, 3, 3, 3, 3, 0}},
		{position: "WR", want: []int{2, 2, 2, 2, 2, 0}},
		{position: "QB", want: []int{1, 1, 1, 1, 1, 0}},
	} {
		if got := skillGroupCurrentLevels(tt.position, groups, ceilings); !slices.Equal(got, tt.want) {
			t.Errorf("%s levels = %v, want %v", tt.position, got, tt.want)
		}
	}

	groups[2].Attributes[0].Rating = nil
	if got := skillGroupCurrentLevels("QB", groups, ceilings); got != nil {
		t.Fatalf("incomplete ratings returned %v, want nil", got)
	}
	if got := skillGroupCurrentLevels("QB", groups[:5], ceilings); got != nil {
		t.Fatalf("five groups returned %v, want nil", got)
	}
}

func skillGroupCapField(i int) string {
	return "SkillGroupCap" + string(rune('0'+i))
}

func TestExportRecruitSkillGroupCaps(t *testing.T) {
	skipIfShortIntegration(t)
	file := openTestSave(t)
	export, err := file.ExportWithOptions(ExportOptions{Sections: ExportSections{Recruits: true}})
	if err != nil {
		t.Fatal(err)
	}

	withCaps := 0
	for _, recruit := range export.Recruits {
		p := recruit.Player
		if p == nil || len(p.SkillGroupCaps) == 0 {
			continue
		}
		withCaps++
		if len(p.SkillGroupCaps) != skillGroupCapCount {
			t.Fatalf("recruit %d: %d caps, want %d", recruit.ID, len(p.SkillGroupCaps), skillGroupCapCount)
		}
		if len(p.SkillGroupUnlockedSlots) != skillGroupCapCount {
			t.Fatalf("recruit %d: %d unlocked slots, want %d", recruit.ID, len(p.SkillGroupUnlockedSlots), skillGroupCapCount)
		}
		calculable := len(p.SkillGroups) == skillGroupCapCount && len(p.SkillGroupUnlockedSlots) == skillGroupCapCount
		for _, group := range p.SkillGroups {
			if _, ok := SkillGroupOverall(group.Attributes); !ok {
				calculable = false
				break
			}
		}
		if calculable && len(p.SkillGroupCurrentLevels) != skillGroupCapCount {
			t.Fatalf("recruit %d: %d current levels, want %d", recruit.ID, len(p.SkillGroupCurrentLevels), skillGroupCapCount)
		}
		if !calculable && len(p.SkillGroupCurrentLevels) != 0 {
			t.Fatalf("recruit %d: current levels emitted for incomplete groups", recruit.ID)
		}
		for i, level := range p.SkillGroupCurrentLevels {
			if level < 0 || level > p.SkillGroupUnlockedSlots[i] {
				t.Fatalf("recruit %d bucket %d: current level %d outside 0..%d", recruit.ID, i+1, level, p.SkillGroupUnlockedSlots[i])
			}
		}
		sumCapped := 0
		sumUnlocked := 0
		for i, capped := range p.SkillGroupCaps {
			sumCapped += capped
			sumUnlocked += p.SkillGroupUnlockedSlots[i]
			if capped+p.SkillGroupUnlockedSlots[i] != skillGroupBucketMax {
				t.Fatalf("recruit %d bucket %d: capped+unlocked = %d, want %d", recruit.ID, i+1, capped+p.SkillGroupUnlockedSlots[i], skillGroupBucketMax)
			}
		}
		if sumCapped != p.SkillGroupCapTotal {
			t.Fatalf("recruit %d: capped total %d != sum %d", recruit.ID, p.SkillGroupCapTotal, sumCapped)
		}
		if sumUnlocked != p.SkillGroupUnlockedTotal {
			t.Fatalf("recruit %d: unlocked total %d != sum %d", recruit.ID, p.SkillGroupUnlockedTotal, sumUnlocked)
		}
	}

	if withCaps == 0 {
		t.Fatal("expected recruits with skill group caps")
	}
	t.Logf("recruits with caps=%d", withCaps)
}

func TestMalachiSingletonSkillGroupCaps(t *testing.T) {
	skipIfShortIntegration(t)
	file := openTestSave(t)
	export, err := file.ExportWithOptions(ExportOptions{Sections: ExportSections{Rosters: true}})
	if err != nil {
		t.Fatal(err)
	}

	wantLabels := []string{"Accuracy", "Power", "IQ", "Elusiveness", "Quickness", "Health"}
	wantCapped := []int{3, 1, 5, 2, 6, 4}
	wantUnlocked := []int{17, 19, 15, 18, 14, 16}

	var player *PlayerExport
	for _, roster := range export.Rosters {
		for i := range roster.Players {
			p := roster.Players[i]
			if p.FirstName == "Malachi" && p.LastName == "Singleton" {
				player = &p
				break
			}
		}
	}
	if player == nil {
		t.Fatal("Malachi Singleton not found")
	}
	if len(player.SkillGroupCurrentLevels) != skillGroupCapCount {
		t.Fatalf("current levels len = %d, want %d", len(player.SkillGroupCurrentLevels), skillGroupCapCount)
	}
	for i, level := range player.SkillGroupCurrentLevels {
		if level < 0 || level > player.SkillGroupUnlockedSlots[i] {
			t.Fatalf("bucket %d: current level %d outside 0..%d", i+1, level, player.SkillGroupUnlockedSlots[i])
		}
	}
	if player.SkillGroupCapTotal != 21 {
		t.Fatalf("skillGroupCapTotal = %d, want 21", player.SkillGroupCapTotal)
	}
	if player.SkillGroupUnlockedTotal != 99 {
		t.Fatalf("skillGroupUnlockedTotal = %d, want 99", player.SkillGroupUnlockedTotal)
	}
	for i := range wantCapped {
		if player.SkillGroupCaps[i] != wantCapped[i] {
			t.Fatalf("skillGroupCaps[%d] = %d, want %d", i, player.SkillGroupCaps[i], wantCapped[i])
		}
		if player.SkillGroupUnlockedSlots[i] != wantUnlocked[i] {
			t.Fatalf("skillGroupUnlockedSlots[%d] = %d, want %d", i, player.SkillGroupUnlockedSlots[i], wantUnlocked[i])
		}
	}
	if len(player.SkillGroups) != skillGroupCapCount {
		t.Fatalf("skillGroups len = %d, want %d", len(player.SkillGroups), skillGroupCapCount)
	}
	for i, group := range player.SkillGroups {
		if group.Label != wantLabels[i] {
			t.Fatalf("skillGroups[%d].label = %q, want %q", i, group.Label, wantLabels[i])
		}
		if group.CappedSlots != wantCapped[i] {
			t.Fatalf("skillGroups[%d].cappedSlots = %d, want %d", i, group.CappedSlots, wantCapped[i])
		}
		if group.UnlockedSlots != wantUnlocked[i] {
			t.Fatalf("skillGroups[%d].unlockedSlots = %d, want %d", i, group.UnlockedSlots, wantUnlocked[i])
		}
	}
}
