package dynasty

import "testing"

func TestRecruitPlayerLinkage(t *testing.T) {
	skipIfShortIntegration(t)
	file := openTestSave(t)
	recruitTable, ok := file.PrimaryTableByName("Recruit")
	if !ok {
		t.Fatal("Recruit table not found")
	}
	if err := recruitTable.ReadRecords(); err != nil {
		t.Fatal(err)
	}

	record := recruitTable.Records[1]
	rank, ok := intFieldOK(record, "NationalRank")
	if !ok || rank <= 0 {
		t.Fatalf("expected national rank on recruit 1, got %d ok=%v", rank, ok)
	}

	playerRef, ok := record.Get("Player")
	if !ok || playerRef.Reference == nil {
		t.Fatal("expected player reference on recruit 1")
	}
	player, _, ok := file.playerRecordByReference(playerRef.Reference)
	if !ok {
		t.Fatal("expected player resolution")
	}
	if stringField(player, "LastName") == "Aarons" {
		t.Fatal("recruit 1 should not resolve to Omar Aarons")
	}
	if stringField(player, "FirstName") == "" {
		t.Fatal("expected player name")
	}

	export, err := file.buildRecruitExports()
	if err != nil {
		t.Fatal(err)
	}
	if len(export) != int(recruitTable.ActiveRecordCount()) {
		t.Fatalf("exports = %d active = %d", len(export), recruitTable.ActiveRecordCount())
	}
	if export[1].NationalRank == nil || *export[1].NationalRank <= 0 {
		t.Fatal("expected national rank in export")
	}
	if export[1].Player == nil || export[1].Player.LastName == "" {
		t.Fatal("expected linked player in export")
	}
}
