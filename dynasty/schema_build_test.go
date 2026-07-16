package dynasty

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectRevisionFromFTX(t *testing.T) {
	src := filepath.Join("..", "data", "cfb27-db-data", "2")
	if _, err := os.Stat(src); err != nil {
		t.Skip("local cfb27-db-data extract not present")
	}
	rev, count, err := detectRevision(src)
	if err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("expected dataRevisionVersion hits")
	}
	if rev != 2 {
		t.Fatalf("revision = %d, want 2", rev)
	}
}

func TestBuildSchemaBundleFromFTX(t *testing.T) {
	src := filepath.Join("..", "data", "cfb27-db-data", "2")
	if _, err := os.Stat(src); err != nil {
		t.Skip("local cfb27-db-data extract not present")
	}
	refPath := filepath.Join("..", "data", "C27_468_2.gz")
	if _, err := os.Stat(refPath); err != nil {
		t.Skip("reference C27_468_2.gz not present")
	}

	out := t.TempDir()
	result, err := BuildSchemaBundle(SchemaBuildOptions{
		Source: src,
		OutDir: out,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Major != 468 || result.Minor != 2 || result.GameYear != 27 {
		t.Fatalf("meta = %d.%d year=%d, want 468.2 year=27", result.Major, result.Minor, result.GameYear)
	}
	if result.MinorSource != "dataRevisionVersion" {
		t.Fatalf("minor source = %q, want dataRevisionVersion", result.MinorSource)
	}
	if result.MajorSource != "default-cfb27" {
		t.Fatalf("major source = %q, want default-cfb27", result.MajorSource)
	}

	built, err := LoadSchemaFile(result.Path)
	if err != nil {
		t.Fatal(err)
	}
	ref, err := LoadSchemaFile(refPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(built.Tables) < len(ref.Tables)-5 {
		t.Fatalf("built tables = %d, ref = %d (too few)", len(built.Tables), len(ref.Tables))
	}

	bp, ok := built.Table("Player")
	if !ok {
		t.Fatal("built missing Player")
	}
	rp, ok := ref.Table("Player")
	if !ok {
		t.Fatal("ref missing Player")
	}
	if len(bp.Attributes) != len(rp.Attributes) {
		t.Fatalf("Player attrs built=%d ref=%d", len(bp.Attributes), len(rp.Attributes))
	}
	builtNames := map[string]string{}
	for _, a := range bp.Attributes {
		builtNames[a.Name] = a.Type
	}
	for _, a := range rp.Attributes {
		if got, ok := builtNames[a.Name]; !ok {
			t.Fatalf("built Player missing field %s", a.Name)
		} else if got != a.Type {
			t.Fatalf("Player.%s type = %s, want %s", a.Name, got, a.Type)
		}
	}
	// Enum embedding
	var pos *FieldSchema
	for i := range bp.Attributes {
		if bp.Attributes[i].Name == "IronManPosition" {
			pos = &bp.Attributes[i]
			break
		}
	}
	if pos == nil || pos.Enum == nil || pos.Enum.Name != "PositionE" {
		t.Fatalf("IronManPosition enum = %+v", pos)
	}
	if len(pos.Enum.Members) < 20 {
		t.Fatalf("PositionE members = %d, want many", len(pos.Enum.Members))
	}
}

func TestResolveSchemaSourcePicksNewestRevision(t *testing.T) {
	src := filepath.Join("..", "data", "cfb27-db-data")
	if _, err := os.Stat(src); err != nil {
		t.Skip("local cfb27-db-data extract not present")
	}
	got, err := resolveSchemaSource(src)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(got) != "2" && filepath.Base(got) != "0" {
		// 0 and 2 are identical rev=2; picker prefers higher dir num when rev ties → 2
		t.Fatalf("picked %q, want patch 2 (or 0)", got)
	}
	if filepath.Base(got) != "2" {
		t.Fatalf("picked %q, want 2 when revisions tie", got)
	}
}
