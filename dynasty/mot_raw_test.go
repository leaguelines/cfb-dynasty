package dynasty

import (
	"fmt"
	"testing"
)

func TestDumpMotivationRaw(t *testing.T) {
	file := openTestSave(t)
	if err := file.ReadAllRecords(); err != nil {
		t.Fatal(err)
	}
	playerTable, _ := file.PrimaryTableByName("Player")
	nonzero := 0
	for _, rec := range playerTable.Records {
		for _, field := range []string{"Motivation1", "Motivation2", "Motivation3"} {
			v, ok := rec.Get(field)
			if !ok {
				continue
			}
			if v.String != "" && v.String != "None" {
				fmt.Printf("player %d %s String=%q Int=%d\n", rec.Index, field, v.String, v.Int)
				nonzero++
			} else if v.Int > 0 {
				fmt.Printf("player %d %s String=%q Int=%d\n", rec.Index, field, v.String, v.Int)
				nonzero++
			}
		}
	}
	fmt.Println("nonzero count", nonzero)
}
