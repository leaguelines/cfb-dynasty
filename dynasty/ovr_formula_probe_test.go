package dynasty

import (
	"fmt"
	"path/filepath"
	"sort"
	"testing"
)

func TestProbeOverallPercentageAndSkillGroupRating(t *testing.T) {
	settings := DefaultSettings()
	settings.SchemaDir = filepath.Join("..", "data")
	settings.TuningPath = filepath.Join("..", "data", "cfb27-db-data", "4", "dynasty-tuning-binary.FTC")
	settings.AutoParse = true

	// Try tuning file tables
	tf, err := Open(settings.TuningPath, &settings)
	if err != nil {
		t.Fatal(err)
	}
	names := []string{}
	for _, n := range []string{"OverallPercentage", "PlayerGradeEval", "PlayerSkillGroup", "PlayerSkillGroupBucket", "PlayerSkill"} {
		if tbl, ok := tf.PrimaryTableByName(n); ok {
			tbl.ReadRecords()
			names = append(names, fmt.Sprintf("%s:%d", n, len(tbl.Records)))
		} else {
			names = append(names, n+":missing")
		}
	}
	fmt.Println("tuning tables:", names)

	if tbl, ok := tf.PrimaryTableByName("OverallPercentage"); ok {
		for i, rec := range tbl.Records {
			if i > 5 {
				break
			}
			keys := []string{}
			for k := range rec.Fields {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			fmt.Printf("OverallPercentage[%d] keys=%v pos=%q\n", i, keys, stringField(rec, "PlayerPosition"))
		}
	}

	// Fit skill group ratings for Malachi vs primary/secondary formulas
	sf := openTestSave(t)
	ex, err := sf.ExportWithOptions(ExportOptions{Sections: ExportSections{Rosters: true}})
	if err != nil {
		t.Fatal(err)
	}
	// game UI from screenshot (approx): Acc 81?, Power 86, IQ 78, Elu 65, Quick 78, Health 86?
	// We trust Power=86 IQ=78 Elu=65 Quick=78 from primary avg match
	for _, r := range ex.Rosters {
		for _, p := range r.Players {
			if p.FirstName != "Malachi" || p.LastName != "Singleton" {
				continue
			}
			fmt.Printf("Malachi OVR=%v archetype=%s/%s\n", p.Overall, p.Archetype, p.ArchetypeLabel)
			for _, g := range p.SkillGroups {
				var prim, sec, ter, all []int
				for _, a := range g.Attributes {
					if a.Rating == nil {
						continue
					}
					all = append(all, *a.Rating)
					switch a.Tier {
					case "Primary":
						prim = append(prim, *a.Rating)
					case "Secondary":
						sec = append(sec, *a.Rating)
					case "Tertiary":
						ter = append(ter, *a.Rating)
					}
				}
				avg := func(xs []int) int {
					if len(xs) == 0 {
						return -1
					}
					s := 0
					for _, v := range xs {
						s += v
					}
					return s / len(xs)
				}
				// weighted guesses
				wSum, wDen := 0, 0
				for _, v := range prim {
					wSum += v * 3
					wDen += 3
				}
				for _, v := range sec {
					wSum += v * 2
					wDen += 2
				}
				for _, v := range ter {
					wSum += v * 1
					wDen += 1
				}
				wAvg := -1
				if wDen > 0 {
					wAvg = wSum / wDen
				}
				fmt.Printf("%s prim=%v(%d) sec=%v(%d) ter=%v(%d) all=%d w3/2/1=%d\n",
					g.Label, prim, avg(prim), sec, avg(sec), ter, avg(ter), avg(all), wAvg)
			}
			return
		}
	}
}
