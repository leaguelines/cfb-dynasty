package dynasty

import (
	"os"
	"path/filepath"
	"testing"
)

func TestArchetypeLabelFromRecord_PTFlags(t *testing.T) {
	tests := []struct {
		field string
		want  string
	}{
		{"PT_QBPOCKETPASSER", "Pocket Passer"},
		{"PT_QBBACKFIELDCREATOR", "Backfield Creator"},
		{"PT_QBDUALTHREAT", "Dual Threat"},
		{"PT_QBPURERUNNER", "Pure Runner"},
		{"PT_HBPOWERBACK", "Contact Seeker"},
		{"PT_WRELUSIVEROUTERUNNER", "Elusive Route Runner"},
		{"PT_WRPOWERBLOCKING", "Speedster"},
		{"PT_DLPUREPOWER", "Pure Power"},
		{"PT_LBPASSCOVERAGE", "Lurker"},
	}

	for _, tt := range tests {
		record := Record{
			Fields: map[string]FieldValue{
				tt.field: {Bool: true},
			},
		}
		if got := archetypeLabelFromRecord(record); got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.field, got, tt.want)
		}
	}
}

func TestArchetypeLabelFromRecord_PlayerTypeFallback(t *testing.T) {
	record := Record{
		Fields: map[string]FieldValue{
			"PlayerType": {String: "CB_Zone"},
		},
	}
	if got := archetypeLabelFromRecord(record); got != "Zone" {
		t.Fatalf("got %q, want Zone", got)
	}
}

func TestArchetypeLabelFromRecord_PTPreferredOverPlayerType(t *testing.T) {
	record := Record{
		Fields: map[string]FieldValue{
			"PlayerType":          {String: "QB_FieldGeneral"},
			"PT_QBBACKFIELDCREATOR": {Bool: true},
		},
	}
	if got := archetypeLabelFromRecord(record); got != "Backfield Creator" {
		t.Fatalf("got %q, want Backfield Creator", got)
	}
}

func TestArchetypeLabelFromRecord_RealSaveQBs(t *testing.T) {
	savePath := filepath.Join("..", "data", "DYNASTY-TESTSAVE-27")
	if _, err := os.Stat(savePath); err != nil {
		t.Skip("test save not available:", savePath)
	}
	file := openTestSave(t)
	table, ok := file.PrimaryTableByName("Player")
	if !ok {
		t.Fatal("Player table missing")
	}
	if err := table.ReadRecords(); err != nil {
		t.Fatal(err)
	}

	labels := map[string]int{}
	for _, record := range table.Records {
		pos, _ := record.Get("Position")
		if pos.String != "QB" {
			continue
		}
		label := archetypeLabelFromRecord(record)
		if label == "" {
			continue
		}
		labels[label]++
	}
	if labels["Pocket Passer"] == 0 {
		t.Fatalf("expected Pocket Passer QBs, got %#v", labels)
	}
	if labels["Dual Threat"] == 0 {
		t.Fatalf("expected Dual Threat QBs, got %#v", labels)
	}
}
