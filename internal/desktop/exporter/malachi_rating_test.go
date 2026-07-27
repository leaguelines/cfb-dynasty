package exporter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

func TestSkillGroupDisplayRatingUsesWeightedOverall(t *testing.T) {
	primary, secondary, tertiary := 80, 100, 70
	attributes := []dynasty.SkillGroupAttributeExport{
		{Rating: &primary, Tier: "Primary"},
		{Rating: &secondary, Tier: "Secondary"},
		{Rating: &tertiary, Tier: "Tertiary"},
	}

	got, ok := skillGroupDisplayRating(attributes, nil)
	if !ok || got != 82 {
		t.Fatalf("skillGroupDisplayRating() = (%d, %v), want (82, true)", got, ok)
	}
}

func TestSkillGroupRowsUseCurrentLevelsForSegments(t *testing.T) {
	player := &dynasty.PlayerExport{
		SkillGroups: []dynasty.SkillGroupExport{
			{Slot: 1, CappedSlots: 3, UnlockedSlots: 17},
			{Slot: 2, CappedSlots: 5, UnlockedSlots: 15},
		},
		SkillGroupCurrentLevels: []int{7, 99},
	}

	rows := SkillGroupRows(player)
	assertSkillSegmentCounts(t, rows[0].Segments, 7, 10, 3)
	assertSkillSegmentCounts(t, rows[1].Segments, 15, 0, 5)

	player.SkillGroupCurrentLevels = nil
	rows = SkillGroupRows(player)
	assertSkillSegmentCounts(t, rows[0].Segments, 0, 17, 3)
}

func TestSkillGroupRowsDefaultMissingFallbackCurrentLevelToZero(t *testing.T) {
	player := &dynasty.PlayerExport{
		SkillGroupCaps:          []int{3, 5},
		SkillGroupUnlockedSlots: []int{17, 15},
		SkillGroupCurrentLevels: []int{7},
	}

	rows := SkillGroupRows(player)
	assertSkillSegmentCounts(t, rows[0].Segments, 7, 10, 3)
	assertSkillSegmentCounts(t, rows[1].Segments, 0, 15, 5)
}

func TestBuildSkillGroupSegmentsNormalizesBounds(t *testing.T) {
	assertSkillSegmentCounts(t, buildSkillGroupSegments(-1, -2), 0, 20, 0)
	assertSkillSegmentCounts(t, buildSkillGroupSegments(99, 25), 0, 0, 20)
}

func TestBuildSkillGroupSegmentsOrdersKinds(t *testing.T) {
	assertSkillSegmentCounts(t, buildSkillGroupSegments(7, 3), 7, 10, 3)
}

func assertSkillSegmentCounts(t *testing.T, segments []SkillGroupSegment, filled, available, capped int) {
	t.Helper()
	counts := map[SkillGroupSegmentKind]int{}
	for i, segment := range segments {
		counts[segment.Kind]++

		want := SkillSegCapped
		if i < filled {
			want = SkillSegFilled
		} else if i < filled+available {
			want = SkillSegAvailable
		}
		if segment.Kind != want {
			t.Fatalf("segment %d = %q, want %q", i, segment.Kind, want)
		}
	}
	if len(segments) != skillGroupBucketMax ||
		counts[SkillSegFilled] != filled ||
		counts[SkillSegAvailable] != available ||
		counts[SkillSegCapped] != capped {
		t.Fatalf("segments = total:%d filled:%d available:%d capped:%d, want total:%d filled:%d available:%d capped:%d",
			len(segments), counts[SkillSegFilled], counts[SkillSegAvailable], counts[SkillSegCapped],
			skillGroupBucketMax, filled, available, capped)
	}
}

func TestMalachiSkillGroupRatingsMatchWeightedOverall(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	savePath := filepath.Join("..", "..", "..", "data", "DYNASTY-TESTSAVE-27")
	if _, err := os.Stat(savePath); err != nil {
		t.Skip("test save not available:", savePath)
	}
	s := dynasty.DefaultSettings()
	s.SchemaDir = filepath.Join("..", "..", "..", "data")
	s.AutoParse = true
	f, err := dynasty.Open(savePath, &s)
	if err != nil {
		t.Fatal(err)
	}
	ex, err := f.ExportWithOptions(dynasty.ExportOptions{Sections: dynasty.ExportSections{Rosters: true}})
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range ex.Rosters {
		for i := range r.Players {
			p := &r.Players[i]
			if p.FirstName != "Malachi" || p.LastName != "Singleton" {
				continue
			}
			rows := SkillGroupRows(p)
			if len(rows) != len(p.SkillGroups) {
				t.Fatalf("rows = %d, want %d", len(rows), len(p.SkillGroups))
			}
			for i, group := range p.SkillGroups {
				for _, attribute := range group.Attributes {
					if attribute.Tier == "" {
						t.Errorf("%s attr %s missing tier", group.Label, attribute.Name)
					}
				}
				want, wantOK := dynasty.SkillGroupOverall(group.Attributes)
				if rows[i].Rating != want || rows[i].HasRating != wantOK {
					t.Errorf("%s rating = (%d, %v), want (%d, %v)", group.Label, rows[i].Rating, rows[i].HasRating, want, wantOK)
				}
			}
			return
		}
	}
	t.Fatal("Malachi Singleton not found")
}
