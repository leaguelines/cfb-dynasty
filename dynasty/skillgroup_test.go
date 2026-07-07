package dynasty

import "testing"

func TestSkillGroupCapsFromRecord(t *testing.T) {
	record := Record{Fields: map[string]FieldValue{
		"SkillGroupCap1": {Int: 20},
		"SkillGroupCap2": {Int: 18},
		"SkillGroupCap3": {Int: 14},
		"SkillGroupCap4": {Int: 12},
		"SkillGroupCap5": {Int: 10},
		"SkillGroupCap6": {Int: 8},
	}}
	caps, total, ok := skillGroupCapsFromRecord(record)
	if !ok {
		t.Fatal("expected ok")
	}
	if total != 82 {
		t.Fatalf("total = %d, want 82", total)
	}
	want := []int{20, 18, 14, 12, 10, 8}
	for i, c := range want {
		if caps[i] != c {
			t.Fatalf("caps[%d] = %d, want %d", i, caps[i], c)
		}
	}
}

func TestSkillGroupCapsFromRecord_MissingOrEmpty(t *testing.T) {
	if _, _, ok := skillGroupCapsFromRecord(Record{Fields: map[string]FieldValue{"SkillGroupCap1": {Int: 5}}}); ok {
		t.Fatal("expected not ok when caps are incomplete")
	}
	zero := Record{Fields: map[string]FieldValue{}}
	for i := 1; i <= skillGroupCapCount; i++ {
		zero.Fields[skillGroupCapField(i)] = FieldValue{Int: 0}
	}
	if _, _, ok := skillGroupCapsFromRecord(zero); ok {
		t.Fatal("expected not ok when all caps are zero")
	}
}

func skillGroupCapField(i int) string {
	return "SkillGroupCap" + string(rune('0'+i))
}

func TestExportRecruitSkillGroupCaps(t *testing.T) {
	file := openTestSave(t)
	export, err := file.Export()
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
		sum := 0
		for _, c := range p.SkillGroupCaps {
			sum += c
		}
		if sum != p.SkillGroupCapTotal {
			t.Fatalf("recruit %d: total %d != sum %d", recruit.ID, p.SkillGroupCapTotal, sum)
		}
	}

	if withCaps == 0 {
		t.Fatal("expected recruits with skill group caps")
	}
	t.Logf("recruits with caps=%d", withCaps)
}
