package dynasty

import (
	"fmt"
	"testing"
)

func TestDumpMotivations(t *testing.T) {
	file := openTestSave(t)
	export, err := file.ExportWithOptions(ExportOptions{Sections: ExportSections{Recruits: true, Rosters: true}})
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]int{}
	for _, r := range export.Recruits {
		if r.Player == nil {
			continue
		}
		for _, m := range r.Player.Motivations {
			seen[m]++
		}
	}
	fmt.Println("recruit motivations:", seen)
	seen2 := map[string]int{}
	for _, roster := range export.Rosters {
		for _, p := range roster.Players {
			for _, m := range p.Motivations {
				seen2[m]++
			}
		}
	}
	top := 0
	for k, v := range seen2 {
		if v > top {
			fmt.Println("roster mot", k, v)
		}
	}
}
