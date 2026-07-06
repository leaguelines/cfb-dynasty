package dynasty

import "testing"

func TestExportSectionsAll(t *testing.T) {
	if !(ExportSections{}).all() {
		t.Fatal("empty sections should mean all")
	}
	if (ExportSections{Games: true}).all() {
		t.Fatal("explicit section should not mean all")
	}
}

func TestExportOptionsIncludeSections(t *testing.T) {
	opts := ExportOptions{Sections: ExportSections{Games: true}}
	if !opts.IncludeGames() {
		t.Fatal("expected games")
	}
	if opts.IncludeRecruits() {
		t.Fatal("did not request recruits")
	}
	if opts.IncludeSeason() {
		t.Fatal("did not request season")
	}
}
