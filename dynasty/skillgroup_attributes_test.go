package dynasty

import (
	"path/filepath"
	"testing"
)

func TestPocketPasserAccuracyAttributes(t *testing.T) {
	schemaDir := filepath.Join("..", "data")
	settings := DefaultSettings()
	settings.SchemaDir = schemaDir
	settings.TuningPath = filepath.Join(schemaDir, "cfb27-db-data", "2", "dynasty-tuning-binary.FTC")
	settings.AutoParse = true
	tf, err := Open(settings.TuningPath, &settings)
	if err != nil {
		t.Skip("tuning data not available:", err)
	}
	idx := buildSkillGroupIndex(tf)
	info := idx.bucketInfo(Record{Fields: map[string]FieldValue{
		"PlayerType":        {String: "QB_FieldGeneral"},
		"PT_QBPOCKETPASSER": {Bool: true},
	}})
	if len(info.attributes) != skillGroupCapCount {
		t.Fatalf("attributes buckets = %d, want %d", len(info.attributes), skillGroupCapCount)
	}
	want := []string{
		"ThrowAccuracyMid",
		"ThrowAccuracyDeep",
		"ThrowAccuracyShort",
		"ThrowUnderPressure",
		"ThrowOnTheRun",
	}
	for i, ability := range want {
		if info.attributes[0][i].playerAbility != ability {
			t.Fatalf("accuracy[%d] = %q, want %q", i, info.attributes[0][i].playerAbility, ability)
		}
	}
	for _, attr := range info.attributes[0] {
		if attr.playerAbility == "KickAccuracy" {
			t.Fatalf("kick accuracy should not be in QB accuracy bucket")
		}
	}
}

func TestMalachiSingletonSkillGroupAttributes(t *testing.T) {
	skipIfShortIntegration(t)
	file := openTestSave(t)
	export, err := file.ExportWithOptions(ExportOptions{Sections: ExportSections{Rosters: true}})
	if err != nil {
		t.Fatal(err)
	}
	for _, roster := range export.Rosters {
		for i := range roster.Players {
			p := roster.Players[i]
			if p.FirstName != "Malachi" || p.LastName != "Singleton" {
				continue
			}
			if len(p.SkillGroups) == 0 {
				t.Fatal("missing skill groups")
			}
			acc := p.SkillGroups[0]
			if acc.Label != "Accuracy" {
				t.Fatalf("first group = %q, want Accuracy", acc.Label)
			}
			if len(acc.Attributes) != 5 {
				t.Fatalf("accuracy attributes = %d, want 5", len(acc.Attributes))
			}
			for _, attr := range acc.Attributes {
				if attr.PlayerAbility == "KickAccuracy" {
					t.Fatal("kick accuracy in accuracy bucket")
				}
				if attr.RatingKey == "" || attr.Rating == nil {
					t.Fatalf("missing rating for %q", attr.Name)
				}
			}
			power := p.SkillGroups[1]
			if power.Label != "Power" {
				t.Fatalf("second group = %q, want Power", power.Label)
			}
			if len(power.Attributes) != 2 {
				t.Fatalf("power attributes = %d, want 2 (%v)", len(power.Attributes), power.Attributes)
			}
			if power.Attributes[0].PlayerAbility != "ThrowPower" || power.Attributes[1].PlayerAbility != "Strength" {
				t.Fatalf("unexpected power attributes: %+v", power.Attributes)
			}
			for _, ability := range []string{"KickPower", "HitPower"} {
				for _, attr := range power.Attributes {
					if attr.PlayerAbility == ability {
						t.Fatalf("%s should not be in power bucket", ability)
					}
				}
			}
			return
		}
	}
	t.Fatal("Malachi Singleton not found")
}

func TestPlayerAbilityRatingKey(t *testing.T) {
	if got := playerAbilityRatingKey("ThrowPower"); got != "ThrowPowerRating" {
		t.Fatalf("got %q", got)
	}
}
