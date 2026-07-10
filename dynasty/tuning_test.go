package dynasty

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSkillGroupIndexFromTuning(t *testing.T) {
	schemaDir := filepath.Join("..", "data")
	tuningPath := filepath.Join(schemaDir, "cfb27-db-data", "2", "dynasty-tuning-binary.FTC")
	if _, err := os.Stat(tuningPath); err != nil {
		t.Skip("tuning data not available:", tuningPath)
	}

	settings := DefaultSettings()
	settings.SchemaDir = schemaDir
	settings.TuningPath = tuningPath
	settings.AutoParse = true
	tf, err := Open(tuningPath, &settings)
	if err != nil {
		t.Fatal(err)
	}
	idx := buildSkillGroupIndex(tf)
	info := idx.bucketInfo(Record{Fields: map[string]FieldValue{
		"PlayerType":        {String: "QB_FieldGeneral"},
		"PT_QBPOCKETPASSER": {Bool: true},
	}})
	if len(info.labels) != skillGroupCapCount {
		t.Fatalf("labels = %#v, want %d entries", info.labels, skillGroupCapCount)
	}
	if len(info.slots) != skillGroupCapCount {
		t.Fatalf("slots = %#v, want %d entries", info.slots, skillGroupCapCount)
	}
	hasAccuracy := false
	for _, name := range info.labels {
		if name == "Accuracy" {
			hasAccuracy = true
			break
		}
	}
	if !hasAccuracy {
		t.Fatalf("expected Accuracy bucket for Pocket Passer, got %#v", info.labels)
	}
	for i, slots := range info.slots {
		if slots <= 0 {
			t.Fatalf("bucket %d (%s): attributeCount = %d, want > 0", i, info.labels[i], slots)
		}
	}
}

func TestApplySkillGroupLabelsOnExport(t *testing.T) {
	skipIfShortIntegration(t)
	file := openTestSave(t)
	export, err := file.ExportWithOptions(ExportOptions{Sections: ExportSections{Rosters: true}})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, roster := range export.Rosters {
		for _, player := range roster.Players {
			if len(player.SkillGroups) == 0 {
				continue
			}
			found = true
			if len(player.SkillGroups) != len(player.SkillGroupCaps) {
				t.Fatalf("player %d: %d groups vs %d caps", player.ID, len(player.SkillGroups), len(player.SkillGroupCaps))
			}
			if player.SkillGroups[0].Label == "" {
				t.Fatalf("player %d: missing skill group label", player.ID)
			}
			if player.SkillGroups[0].AttributeCount <= 0 {
				t.Fatalf("player %d: missing skill group attribute count", player.ID)
			}
			if len(player.SkillGroupAttributeCounts) != len(player.SkillGroupCaps) {
				t.Fatalf("player %d: %d attribute counts vs %d capped slots", player.ID, len(player.SkillGroupAttributeCounts), len(player.SkillGroupCaps))
			}
			if player.SkillGroupCaps[0]+player.SkillGroupUnlockedSlots[0] != skillGroupBucketMax {
				t.Fatalf("player %d: capped+unlocked != %d", player.ID, skillGroupBucketMax)
			}
			return
		}
	}
	if !found {
		t.Skip("no labeled skill groups in test save (tuning path may be missing)")
	}
}
