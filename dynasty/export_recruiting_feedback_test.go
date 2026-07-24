package dynasty

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRecruitingFeedbackHardSellSlotReload(t *testing.T) {
	skipIfShortIntegration(t)
	dir := filepath.Join("..", "data", "hard-sell-slots-test")
	if _, err := os.Stat(dir); err != nil {
		t.Skip("hard-sell-slots-test saves not available:", dir)
	}
	s := DefaultSettings()
	s.SchemaDir = filepath.Join("..", "data")
	s.AutoParse = true

	type result struct {
		influence int
		gained    int
		pipeline  int
	}
	want := map[string]result{
		"DYNASTY-BASELINE": {influence: 126},
		"DYNASTY-S1":       {influence: 190, gained: 64, pipeline: 20},
		"DYNASTY-S2":       {influence: 194, gained: 68, pipeline: 20},
		"DYNASTY-S3":       {influence: 190, gained: 64, pipeline: 20},
	}

	for name, expect := range want {
		path := filepath.Join(dir, name)
		sf, err := Open(path, &s)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		ex, err := sf.ExportWithOptions(ExportOptions{
			Sections: ExportSections{Recruits: true, Recruiting: true},
		})
		if err != nil {
			t.Fatalf("%s export: %v", name, err)
		}
		var target *RecruitingTargetExport
		for i := range ex.Recruiting {
			if ex.Recruiting[i].RecruitID == 2315 {
				target = &ex.Recruiting[i]
				break
			}
		}
		if target == nil {
			t.Fatalf("%s: missing recruit 2315 target", name)
		}
		ru := 0
		for _, si := range target.SchoolInterest {
			if si.TeamName == "Rutgers" {
				ru = si.Influence
			}
		}
		if ru != expect.influence {
			t.Fatalf("%s Rutgers influence=%d want %d", name, ru, expect.influence)
		}
		if expect.gained == 0 {
			continue
		}
		if len(target.RecruitingFeedback) == 0 {
			t.Fatalf("%s: expected RecruitingFeedback", name)
		}
		fb := target.RecruitingFeedback[0]
		if fb.Intensity != "HardSell" {
			t.Fatalf("%s intensity=%q", name, fb.Intensity)
		}
		if fb.InfluenceGained == nil || *fb.InfluenceGained != expect.gained {
			t.Fatalf("%s gained=%v want %d", name, fb.InfluenceGained, expect.gained)
		}
		if fb.MaxInfluenceGain == nil || *fb.MaxInfluenceGain != expect.gained {
			t.Fatalf("%s max=%v want %d (max mirrors rolled total, not perm band)", name, fb.MaxInfluenceGain, expect.gained)
		}
		pipe := 0
		for _, b := range fb.Bonuses {
			if b.BonusType == "Pipeline" {
				pipe = b.BonusValue
			}
		}
		if pipe != expect.pipeline {
			t.Fatalf("%s pipeline bonus=%d want %d", name, pipe, expect.pipeline)
		}
		gradePart := *fb.InfluenceGained - pipe
		if gradePart != 44 && gradePart != 48 {
			t.Fatalf("%s grade+base=%d want 44 (FTC) or 48 (alt perm)", name, gradePart)
		}
	}
}
