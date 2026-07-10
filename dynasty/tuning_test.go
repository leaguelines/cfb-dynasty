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
	buckets := idx.bucketLabels(Record{Fields: map[string]FieldValue{
		"PlayerType":        {String: "QB_FieldGeneral"},
		"PT_QBPOCKETPASSER": {Bool: true},
	}})
	if len(buckets) != skillGroupCapCount {
		t.Fatalf("buckets = %#v, want %d entries", buckets, skillGroupCapCount)
	}
	hasAccuracy := false
	for _, name := range buckets {
		if name == "Accuracy" {
			hasAccuracy = true
			break
		}
	}
	if !hasAccuracy {
		t.Fatalf("expected Accuracy bucket for Pocket Passer, got %#v", buckets)
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
			return
		}
	}
	if !found {
		t.Skip("no labeled skill groups in test save (tuning path may be missing)")
	}
}
