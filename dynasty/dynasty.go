package dynasty

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/leaguelines/cfb-dynasty/internal/binary"
	"github.com/leaguelines/cfb-dynasty/internal/compress"
)

// FileInfo summarizes container-level details without full parsing.
type FileInfo struct {
	Path              string
	Size              int64
	SHA256            string
	Compressed        bool
	Format            Format
	GameYear          int
	ExpectedSchema    *SchemaVersion
	HeaderHex         string
	UnpackedSize      int
	TableMarkerCounts map[TableMarker]int
}

// File represents an opened dynasty save on disk.
type File struct {
	path           string
	settings       Settings
	rawContents    []byte
	unpacked       []byte
	info           FileInfo
	loaded         bool
	tables         []Table
	schema         *Schema
	assetTable     []AssetEntry
}

// AssetEntry maps an asset id to a record reference, mirroring Madden's asset table.
type AssetEntry struct {
	AssetID   uint32
	Reference uint32
}

// Open reads a dynasty save from path. Pass nil for default settings.
func Open(path string, settings *Settings) (*File, error) {
	cfg := DefaultSettings()
	if settings != nil {
		cfg = *settings
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cfb-dynasty: read save: %w", err)
	}

	f := &File{
		path:        path,
		settings:    cfg,
		rawContents: raw,
	}

	if err := f.prepare(); err != nil {
		return nil, err
	}

	if cfg.AutoParse {
		if err := f.Parse(); err != nil {
			return nil, err
		}
	}

	return f, nil
}

// Path returns the filesystem path used to open the save.
func (f *File) Path() string {
	return f.path
}

// Settings returns a copy of the active settings.
func (f *File) Settings() Settings {
	return f.settings
}

// RawContents returns the on-disk bytes.
func (f *File) RawContents() []byte {
	return f.rawContents
}

// UnpackedContents returns decompressed bytes when available, otherwise raw data.
func (f *File) UnpackedContents() []byte {
	if len(f.unpacked) > 0 {
		return f.unpacked
	}
	return f.rawContents
}

// Loaded reports whether Parse completed successfully.
func (f *File) Loaded() bool {
	return f.loaded
}

// Tables returns parsed tables after Parse.
func (f *File) Tables() []Table {
	return f.tables
}

// Schema returns loaded schema metadata after Parse.
func (f *File) Schema() *Schema {
	return f.schema
}

// Inspect returns container metadata without requiring Parse.
func (f *File) Inspect() (FileInfo, error) {
	return f.info, nil
}

// Parse unpacks the save, discovers tables, and loads schemas.
func (f *File) Parse() error {
	if f.loaded {
		return nil
	}

	if err := f.prepare(); err != nil {
		return err
	}

	// Table discovery and schema loading depend on confirmed CFB layout.
	if f.settings.SchemaDir != "" {
		if err := f.loadSchema(); err != nil {
			return err
		}
	}

	tables, err := discoverTables(f.unpacked)
	if err != nil {
		return err
	}
	attachTableSchemas(tables, f.schema)
	f.tables = tables
	f.loaded = true
	return nil
}

// Export builds the normalized export payload. Requires Parse to succeed.
func (f *File) Export() (Export, error) {
	return f.ExportWithOptions(DefaultExportOptions())
}

// ExportWithOptions builds an export payload including only the requested sections.
func (f *File) ExportWithOptions(opts ExportOptions) (Export, error) {
	if !f.loaded {
		return Export{}, ErrNotLoaded
	}
	return f.buildExport(opts)
}

// GetTableByName returns the first table with the given name.
func (f *File) GetTableByName(name string) (*Table, bool) {
	for i := range f.tables {
		if f.tables[i].Name() == name {
			return &f.tables[i], true
		}
	}
	return nil, false
}

// GetTableByID returns a table by header table id.
func (f *File) GetTableByID(id uint32) (*Table, bool) {
	for i := range f.tables {
		if f.tables[i].Header.TableID == id {
			return &f.tables[i], true
		}
	}
	return nil, false
}

func (f *File) prepare() error {
	if len(f.rawContents) == 0 {
		return fmt.Errorf("%w: empty file", ErrInvalidSave)
	}

	detected := DetectFormat(f.rawContents)
	unpacked := f.rawContents
	if detected.Compressed {
		data, err := compress.Unpack(f.rawContents, compressContainer(detected))
		if err != nil {
			return fmt.Errorf("cfb-dynasty: unpack: %w", err)
		}
		unpacked = data
		// Re-detect on inflated payload when the compressed wrapper hides FrTk.
		if inner := DetectFormat(unpacked); inner.Format == FormatUncompressed {
			detected = inner
		}
	}

	sum := sha256.Sum256(f.rawContents)
	headerLen := 64
	if len(f.rawContents) < headerLen {
		headerLen = len(f.rawContents)
	}

	f.unpacked = unpacked
	f.info = FileInfo{
		Path:              f.path,
		Size:              int64(len(f.rawContents)),
		SHA256:            hex.EncodeToString(sum[:]),
		Compressed:        detected.Compressed,
		Format:            detected.Format,
		GameYear:          f.settings.GameYear(),
		ExpectedSchema:    detected.Schema,
		HeaderHex:         binary.HexPrefix(f.rawContents, headerLen),
		UnpackedSize:      len(unpacked),
		TableMarkerCounts: countTableMarkers(unpacked),
	}
	return nil
}

func (f *File) loadSchema() error {
	if f.settings.SchemaDir == "" {
		return ErrSchemaRequired
	}
	version := SchemaVersion{}
	if f.settings.SchemaOverride != nil {
		version = *f.settings.SchemaOverride
	} else if f.info.ExpectedSchema != nil {
		version = *f.info.ExpectedSchema
	}
	schema, err := LoadSchema(f.settings.SchemaDir, version)
	if err != nil {
		return err
	}
	f.schema = schema
	return nil
}

func countTableMarkers(unpacked []byte) map[TableMarker]int {
	counts := make(map[TableMarker]int, len(AllTableMarkers))
	for _, marker := range AllTableMarkers {
		counts[marker] = len(binary.FindAll(unpacked, []byte(marker)))
	}
	return counts
}

func compressContainer(detected DetectedFormat) compress.Container {
	if !detected.Compressed {
		return compress.Container{Compressed: false}
	}
	skip := 0
	if detected.Format == FormatDynasty {
		skip = CompressedDataOffset
	}
	return compress.Container{Compressed: true, Skip: skip}
}
