package dynasty

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestListRecruitingTuningTables(t *testing.T) {
	schemaDir := filepath.Join("..", "data")
	settings := DefaultSettings()
	settings.SchemaDir = schemaDir
	settings.AutoParse = true
	path := filepath.Join(schemaDir, "cfb27-db-data", "2", "dynasty-tuning-binary.FTC")
	tf, err := Open(path, &settings)
	if err != nil {
		t.Skip("tuning data not available:", err)
	}
	var names []string
	for i := range tf.tables {
		n := tf.tables[i].Name()
		low := strings.ToLower(n)
		if strings.Contains(low, "recruit") || strings.Contains(low, "pitch") || strings.Contains(low, "visit") || strings.Contains(low, "pipeline") {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	for _, n := range names {
		tbl, _ := tf.GetTableByName(n)
		t.Logf("%s allocated=%d", n, tbl.AllocatedRecordCount())
	}
}
