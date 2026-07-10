package dynasty

import "testing"

func TestArrayStoreMemberRefsPreservesOrder(t *testing.T) {
	settings := DefaultSettings()
	settings.SchemaDir = "../data"
	settings.TuningPath = "../data/cfb27-db-data/2/dynasty-tuning-binary.FTC"
	settings.AutoParse = true
	tf, err := Open(settings.TuningPath, &settings)
	if err != nil {
		t.Skip("tuning data not available:", err)
	}

	sg, ok := tf.PrimaryTableByName("PlayerSkillGroup")
	if !ok {
		t.Fatal("PlayerSkillGroup table missing")
	}
	if err := sg.ReadRecords(); err != nil {
		t.Fatal(err)
	}
	for _, rec := range sg.Records {
		if stringField(rec, "Name") != "Pocket Passer" {
			continue
		}
		ref, ok := rec.Get("PlayerSkillsBucket")
		if !ok || ref.Reference == nil {
			t.Fatal("missing PlayerSkillsBucket")
		}
		var names []string
		for _, memberRef := range tf.arrayStoreMemberRefs(ref.Reference) {
			row, ok := tf.RecordByReference("PlayerSkillGroupBucket", memberRef)
			if !ok {
				continue
			}
			name := stringField(row, "Name")
			if name == "" {
				continue
			}
			names = append(names, name)
			if len(names) >= skillGroupCapCount {
				break
			}
		}
		want := []string{"Accuracy", "Power", "IQ", "Elusiveness", "Quickness", "Health"}
		if len(names) != len(want) {
			t.Fatalf("bucket count = %d, want %d (%v)", len(names), len(want), names)
		}
		for i := range want {
			if names[i] != want[i] {
				t.Fatalf("bucket[%d] = %q, want %q (full=%v)", i, names[i], want[i], names)
			}
		}
		return
	}
	t.Fatal("Pocket Passer skill group not found")
}
