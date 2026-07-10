package dynasty

import (
	"fmt"
	"strings"
	"testing"
)

func TestDumpPitchTables(t *testing.T) {
	skipIfShortIntegration(t)
	file := openTestSave(t)
	if err := file.ReadAllRecords(); err != nil {
		t.Fatal(err)
	}
	for i := range file.tables {
		tbl := &file.tables[i]
		name := tbl.Name()
		if !strings.Contains(strings.ToLower(name), "pitch") && !strings.Contains(strings.ToLower(name), "motivation") {
			continue
		}
		fmt.Println("TABLE", name, "records", len(tbl.Records))
		if tbl.Schema != nil {
			for _, a := range tbl.Schema.Attributes {
				fmt.Println("  attr", a.Name)
			}
		}
		if len(tbl.Records) > 0 && len(tbl.Records[0].Fields) > 0 {
			for k := range tbl.Records[0].Fields {
				if strings.Contains(k, "Pitch") || strings.Contains(k, "Motivation") {
					fmt.Println("  field", k)
				}
			}
		}
	}
}
