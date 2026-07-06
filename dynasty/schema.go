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
	Version  SchemaVersion
	GameYear int
	Tables   map[string]*TableSchema
}

// Table returns a table schema by name.
func (s *Schema) Table(name string) (*TableSchema, bool) {
	if s == nil || s.Tables == nil {
		return nil, false
	}
	t, ok := s.Tables[name]
	return t, ok
}

// TableSchema describes columns for a named table.
type TableSchema struct {
	Name       string
	Base       string
	NumMembers int
	AssetID    uint32
	Attributes []FieldSchema
}

// FieldSchema describes one column in a table schema.
type FieldSchema struct {
	Index     int
	Name      string
	Type      string
	MinValue  string
	MaxValue  string
	MaxLength string
	Default   string
	Final     bool
	Const     bool
	Enum      *EnumSchema
}

// LoadSchema reads schema definitions from dir for the given version.
// Schema files are gzip-compressed JSON bundles (for example C27_441_0.gz).
// When version.Path is set, that file is loaded directly. Otherwise dir is
// scanned for .gz files and the closest major/minor match is chosen.
func LoadSchema(dir string, version SchemaVersion) (*Schema, error) {
	path := version.Path
	if path == "" {
		var err error
		path, err = pickSchemaFile(dir, version)
		if err != nil {
			return nil, err
		}
	}
	return LoadSchemaFile(path)
}

// LoadSchemaFile loads a gzip-compressed schema bundle from path.
func LoadSchemaFile(path string) (*Schema, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, err
	}
	return parseSchemaGzip(data, path)
}
