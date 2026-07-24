package exporter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

func TestMalachiPrimarySkillGroupRatings(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	savePath := filepath.Join("..", "..", "..", "data", "DYNASTY-TESTSAVE-27")
	if _, err := os.Stat(savePath); err != nil {
		t.Skip("test save not available:", savePath)
	}
	s := dynasty.DefaultSettings()
	s.SchemaDir = filepath.Join("..", "..", "..", "data")
	s.AutoParse = true
	f, err := dynasty.Open(savePath, &s)
	if err != nil {
		t.Fatal(err)
	}
	ex, err := f.ExportWithOptions(dynasty.ExportOptions{Sections: dynasty.ExportSections{Rosters: true}})
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]int{
		"Accuracy": 80, "Power": 86, "IQ": 78, "Elusiveness": 65, "Quickness": 78, "Health": 85,
	}
	for _, r := range ex.Rosters {
		for i := range r.Players {
			p := &r.Players[i]
			if p.FirstName != "Malachi" || p.LastName != "Singleton" {
				continue
			}
			for _, g := range p.SkillGroups {
				for _, a := range g.Attributes {
					if a.Tier == "" {
						t.Errorf("%s attr %s missing tier", g.Label, a.Name)
					}
				}
			}
			rows := SkillGroupRows(p)
			for _, row := range rows {
				w, ok := want[row.Label]
				if !ok {
					continue
				}
				if row.Rating != w {
					t.Errorf("%s rating=%d want %d", row.Label, row.Rating, w)
				}
			}
			return
		}
	}
	t.Fatal("Malachi Singleton not found")
}
