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
	// Prefer highest dataRevisionVersion (patch 3 when present).
	want := "2"
	if _, err := os.Stat(filepath.Join(src, "3")); err == nil {
		want = "3"
	}
	if filepath.Base(got) != want {
		t.Fatalf("picked %q, want %s", got, want)
	}
}

func TestDetectSchemaMetaPrefersFranchiseRoot(t *testing.T) {
	src := filepath.Join("..", "data", "cfb27-db-data", "3")
	if _, err := os.Stat(filepath.Join(src, "franchise-schemas.FTX")); err != nil {
		t.Skip("patch 3 franchise-schemas.FTX not present")
	}
	meta, majorSrc, minorSrc, err := detectSchemaMeta(src)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Major != 472 {
		t.Fatalf("major = %d, want 472 from franchise-schemas (not Core 55)", meta.Major)
	}
	if meta.Minor != 0 {
		t.Fatalf("minor = %d, want 0 from franchise-schemas dataMinorVersion", meta.Minor)
	}
	if meta.GameYear != 27 {
		t.Fatalf("gameYear = %d, want 27", meta.GameYear)
	}
	if majorSrc != "franchise-dataMajorVersion" {
		t.Fatalf("majorSource = %q", majorSrc)
	}
	if minorSrc != "franchise-dataMinorVersion" {
		t.Fatalf("minorSource = %q", minorSrc)
	}
}

func TestBuildSchemaBundlePatch3UsesFranchiseVersions(t *testing.T) {
	src := filepath.Join("..", "data", "cfb27-db-data", "3")
	if _, err := os.Stat(filepath.Join(src, "franchise-schemas.FTX")); err != nil {
		t.Skip("patch 3 franchise-schemas.FTX not present")
	}
	out := t.TempDir()
	result, err := BuildSchemaBundle(SchemaBuildOptions{Source: src, OutDir: out})
	if err != nil {
		t.Fatal(err)
	}
	if result.Major != 472 || result.Minor != 0 || result.GameYear != 27 {
		t.Fatalf("meta = %d.%d year=%d, want 472.0 year=27", result.Major, result.Minor, result.GameYear)
	}
	if filepath.Base(result.Path) != "C27_472_0.gz" {
		t.Fatalf("output = %s, want C27_472_0.gz", result.Path)
	}
	if _, err := LoadSchemaFile(result.Path); err != nil {
		t.Fatal(err)
	}
}

