package dynasty

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExportRecruitingTunables(t *testing.T) {
	schemaDir := filepath.Join("..", "data")
	if _, err := os.Stat(schemaDir); err != nil {
		t.Skip("schema dir not available")
	}
	export, err := ExportRecruitingTunables(schemaDir, "")
	if err != nil {
		t.Fatal(err)
	}
	if export.Scalars["MaxRecruitingBoardTargets"] == nil {
		t.Fatalf("expected MaxRecruitingBoardTargets scalar, got %#v", export.Scalars)
	}
	odds := export.Arrays["InstantCommitOddsPerStarLevel"]
	if len(odds) == 0 {
		t.Fatalf("expected InstantCommitOddsPerStarLevel array, got %#v", export.Arrays)
	}
	if odds[0] > 100 || odds[0] <= 0 {
		t.Fatalf("InstantCommitOddsPerStarLevel[0] = %d, want plausible percentage", odds[0])
	}
	if len(export.Pitches) < 10 {
		t.Fatalf("expected recruiting pitches, got %d", len(export.Pitches))
	}
	if export.Visit == nil {
		t.Fatal("expected visit tunables")
	}
	if len(export.Actions) == 0 {
		t.Fatal("expected recruiting actions")
	}
	if export.HighSchool == nil || len(export.HighSchool.Scalars) == 0 {
		t.Fatal("expected high school recruiting tunables")
	}
}
