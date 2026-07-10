package dynasty

import (
	"fmt"
	"testing"
)

func TestDumpActivePitches(t *testing.T) {
	skipIfShortIntegration(t)
	file := openTestSave(t)
	if err := file.ReadAllRecords(); err != nil {
		t.Fatal(err)
	}
	tbl, _ := file.PrimaryTableByName("ActiveRecruitingPitch")
	seen := map[string]int{}
	for _, rec := range tbl.Records {
		p := normalizeEnum(stringField(rec, "Pitch"))
		if p != "" {
			seen[p]++
		}
	}
	for k, v := range seen {
		fmt.Println(k, v)
	}
}
