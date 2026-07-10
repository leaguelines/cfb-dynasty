package dynasty

import "testing"

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
