package dynasty

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSchemaFromGzipFixture(t *testing.T) {
	dir := writeSchemaFixture(t)

	schema, err := LoadSchema(dir, SchemaVersion{Major: 441, Minor: 0})
	if err != nil {
		t.Fatal(err)
	}
	if schema.Version.Major != 441 || schema.Version.Minor != 0 {
		t.Fatalf("version = %+v, want major=441 minor=0", schema.Version)
	}
	if schema.GameYear != 27 {
		t.Fatalf("gameYear = %d, want 27", schema.GameYear)
	}

	seasonGame, ok := schema.Table("SeasonGame")
	if !ok {
		t.Fatal("SeasonGame schema missing")
	}
	if seasonGame.NumMembers != 2 {
		t.Fatalf("SeasonGame numMembers = %d, want 2", seasonGame.NumMembers)
	}
	if len(seasonGame.Attributes) != 2 {
		t.Fatalf("SeasonGame attributes = %d, want 2", len(seasonGame.Attributes))
	}
	if seasonGame.Attributes[0].Name != "HomeScore" || seasonGame.Attributes[0].Index != 0 {
		t.Fatalf("first attribute = %+v, want HomeScore index 0", seasonGame.Attributes[0])
	}
}

func TestLoadSchemaPicksClosestMajor(t *testing.T) {
	dir := writeSchemaFixture(t)

	// Save reports 809.1; only 441.0 is on disk.
	schema, err := LoadSchema(dir, SchemaVersion{Major: 809, Minor: 1})
	if err != nil {
		t.Fatal(err)
	}
	if schema.Version.Major != 441 {
		t.Fatalf("picked major = %d, want 441", schema.Version.Major)
	}
}

func TestLoadSchemaFileDirect(t *testing.T) {
	dir := writeSchemaFixture(t)
	path := filepath.Join(dir, "C27_441_0.gz")

	schema, err := LoadSchemaFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if schema.Version.Path != path {
		t.Fatalf("path = %q, want %q", schema.Version.Path, path)
	}
}

func TestLoadSchemaNotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadSchema(dir, SchemaVersion{Major: 1, Minor: 0})
	if err == nil {
		t.Fatal("expected error for empty schema dir")
	}
	if !errors.Is(err, ErrSchemaNotFound) {
		t.Fatalf("err = %v, want ErrSchemaNotFound", err)
	}
}

func TestLoadSchemaRealBundle(t *testing.T) {
	path := filepath.Join("..", "data", "C27_441_0.gz")
	if _, err := os.Stat(path); err != nil {
		t.Skip("real schema bundle not available:", path)
	}

	schema, err := LoadSchemaFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(schema.Tables) < 3000 {
		t.Fatalf("tables = %d, want at least 3000", len(schema.Tables))
	}

	seasonGame, ok := schema.Table("SeasonGame")
	if !ok {
		t.Fatal("SeasonGame missing")
	}
	if seasonGame.NumMembers != 69 {
		t.Fatalf("SeasonGame numMembers = %d, want 69", seasonGame.NumMembers)
	}

	var homeTeam *FieldSchema
	for i := range seasonGame.Attributes {
		if seasonGame.Attributes[i].Name == "HomeTeam" {
			homeTeam = &seasonGame.Attributes[i]
			break
		}
	}
	if homeTeam == nil {
		t.Fatal("HomeTeam field missing")
	}
	if homeTeam.Type != "Team" || homeTeam.Index != 34 {
		t.Fatalf("HomeTeam = %+v, want type Team index 34", *homeTeam)
	}
}

func writeSchemaFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "C27_441_0.gz")

	bundle := map[string]any{
		"meta": map[string]any{
			"major":    441,
			"minor":    0,
			"gameYear": 27,
		},
		"schemas": []map[string]any{
			{
				"name":       "SeasonGame",
				"numMembers": "2",
				"base":       "FootballSeasonGame",
				"attributes": []map[string]any{
					{"index": "0", "name": "HomeScore", "type": "int"},
					{"index": "1", "name": "AwayScore", "type": "int"},
				},
			},
		},
		"schemaMap": map[string]any{
			"SeasonGame": map[string]any{
				"assetId":    12345,
				"name":       "SeasonGame",
				"numMembers": "2",
				"base":       "FootballSeasonGame",
				"attributes": []map[string]any{
					{"index": "0", "name": "HomeScore", "type": "int", "minValue": "0", "maxValue": "255"},
					{"index": "1", "name": "AwayScore", "type": "int"},
				},
			},
		},
	}

	raw, err := json.Marshal(bundle)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	if _, err := gw.Write(raw); err != nil {
		t.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}
