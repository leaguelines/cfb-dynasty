package dynasty

// SchemaVersion identifies a FranTk-style schema bundle embedded in or paired
// with a save file.
type SchemaVersion struct {
	Major int    `json:"major"`
	Minor int    `json:"minor"`
	Path  string `json:"path,omitempty"`
}

// Schema holds table and field definitions loaded from game assets.
type Schema struct {
	Version SchemaVersion
	Tables  map[string]*TableSchema
}

// TableSchema describes columns for a named table.
type TableSchema struct {
	Name       string
	Base       string
	Attributes []FieldSchema
}

// FieldSchema describes one column in a table schema.
type FieldSchema struct {
	Name string
	Type string
}

// LoadSchema reads schema definitions from dir for the given version.
// Returns ErrNotImplemented until CFB schema assets are extracted.
func LoadSchema(dir string, version SchemaVersion) (*Schema, error) {
	_ = dir
	_ = version
	return nil, ErrNotImplemented
}
