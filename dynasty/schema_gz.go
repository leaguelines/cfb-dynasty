package dynasty

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type schemaBundle struct {
	Meta struct {
		Major    int `json:"major"`
		Minor    int `json:"minor"`
		GameYear int `json:"gameYear"`
	} `json:"meta"`
	Schemas   []rawTableSchema            `json:"schemas"`
	SchemaMap map[string]rawTableSchema   `json:"schemaMap"`
}

type rawTableSchema struct {
	AssetID    jsonUint32           `json:"assetId"`
	Name       string             `json:"name"`
	Base       string             `json:"base"`
	NumMembers jsonStringInt      `json:"numMembers"`
	Attributes []rawFieldSchema   `json:"attributes"`
}

type rawFieldSchema struct {
	Index     jsonStringInt `json:"index"`
	Name      string        `json:"name"`
	Type      string        `json:"type"`
	MinValue  string        `json:"minValue"`
	MaxValue  string        `json:"maxValue"`
	MaxLength string        `json:"maxLength"`
	Default   string        `json:"default"`
	Final     string        `json:"final"`
	Const     string        `json:"const"`
	Enum      *rawEnum      `json:"enum"`
}

type rawEnum struct {
	Name      string          `json:"_name"`
	MaxLength jsonStringInt   `json:"_maxLength"`
	Members   []rawEnumMember `json:"_members"`
}

type rawEnumMember struct {
	Name              string        `json:"_name"`
	Value             jsonStringInt `json:"_value"`
	UnformattedValue  string        `json:"_unformattedValue"`
}

// jsonStringInt unmarshals schema index fields that may be JSON strings or numbers.
type jsonStringInt int

// jsonUint32 unmarshals schema asset IDs that may be JSON strings or numbers.
type jsonUint32 uint32

func (v *jsonUint32) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*v = 0
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		n, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return err
		}
		*v = jsonUint32(n)
		return nil
	}
	var n uint32
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*v = jsonUint32(n)
	return nil
}

func (v *jsonStringInt) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*v = 0
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*v = jsonStringInt(n)
		return nil
	}
	var n int
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*v = jsonStringInt(n)
	return nil
}

func readFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cfb-dynasty: read schema %q: %w", path, err)
	}
	return data, nil
}

func parseSchemaGzip(data []byte, path string) (*Schema, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("cfb-dynasty: schema gzip %q: %w", path, err)
	}
	defer reader.Close()

	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("cfb-dynasty: schema inflate %q: %w", path, err)
	}

	var bundle schemaBundle
	if err := json.Unmarshal(raw, &bundle); err != nil {
		return nil, fmt.Errorf("cfb-dynasty: schema json %q: %w", path, err)
	}

	tables := make(map[string]*TableSchema)
	for i := range bundle.Schemas {
		rawTable := bundle.Schemas[i]
		if rawTable.Name == "" {
			continue
		}
		tables[rawTable.Name] = convertTableSchema(rawTable)
	}
	for name, rawTable := range bundle.SchemaMap {
		if rawTable.Name == "" {
			rawTable.Name = name
		}
		tables[rawTable.Name] = convertTableSchema(rawTable)
	}

	if len(tables) == 0 {
		return nil, fmt.Errorf("%w: no tables in %q", ErrInvalidSchema, path)
	}

	version := SchemaVersion{
		Major: bundle.Meta.Major,
		Minor: bundle.Meta.Minor,
		Path:  path,
	}

	return &Schema{
		Version:  version,
		GameYear: bundle.Meta.GameYear,
		Tables:   tables,
	}, nil
}

func convertTableSchema(raw rawTableSchema) *TableSchema {
	attrs := make([]FieldSchema, 0, len(raw.Attributes))
	for _, field := range raw.Attributes {
		attrs = append(attrs, FieldSchema{
			Index:     int(field.Index),
			Name:      field.Name,
			Type:      field.Type,
			MinValue:  field.MinValue,
			MaxValue:  field.MaxValue,
			MaxLength: field.MaxLength,
			Default:   field.Default,
			Final:     field.Final == "true",
			Const:     field.Const == "true",
			Enum:      convertEnum(field.Enum),
		})
	}
	return &TableSchema{
		Name:       raw.Name,
		Base:       raw.Base,
		NumMembers: int(raw.NumMembers),
		AssetID:    uint32(raw.AssetID),
		Attributes: attrs,
	}
}

func convertEnum(raw *rawEnum) *EnumSchema {
	if raw == nil {
		return nil
	}
	enum := &EnumSchema{
		Name:      raw.Name,
		MaxLength: int(raw.MaxLength),
	}
	for _, member := range raw.Members {
		enum.Members = append(enum.Members, EnumMember{
			Name:             member.Name,
			Value:            int(member.Value),
			UnformattedValue: member.UnformattedValue,
		})
	}
	return enum
}
