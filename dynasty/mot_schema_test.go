package dynasty

import (
	"fmt"
	"strings"
	"testing"
)

func TestFindMotivationFields(t *testing.T) {
	file := openTestSave(t)
	if err := file.ReadAllRecords(); err != nil {
		t.Fatal(err)
	}
	recruitTable, _ := file.PrimaryTableByName("Recruit")
	if recruitTable != nil && recruitTable.Schema != nil {
		for _, attr := range recruitTable.Schema.Attributes {
			if strings.Contains(strings.ToLower(attr.Name), "motiv") || strings.Contains(attr.Name, "Pitch") {
				fmt.Println("Recruit", attr.Name)
			}
		}
	}
	// sample dealbreaker values from recruits
	export, _ := file.ExportWithOptions(ExportOptions{Sections: ExportSections{Recruits: true}})
	db := map[string]int{}
	for _, r := range export.Recruits {
		if r.Player != nil && r.Player.RecruitingDealbreaker != "" {
			db[r.Player.RecruitingDealbreaker]++
		}
	}
	fmt.Println("dealbreakers", db)
}
