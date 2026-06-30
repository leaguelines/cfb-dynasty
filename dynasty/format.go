package dynasty

import (
	"encoding/binary"

	pkgbinary "github.com/leaguelines/cfb-dynasty/internal/binary"
)

// DetectedFormat holds container detection results.
type DetectedFormat struct {
	Compressed bool
	Format     Format
	Schema     *SchemaVersion
}

// DetectFormat inspects raw save bytes and returns the best-known container type.
func DetectFormat(data []byte) DetectedFormat {
	if len(data) < 4 {
		return DetectedFormat{Format: FormatUnknown}
	}

	if pkgbinary.HasPrefix(data, []byte(MagicFrTk)) {
		return DetectedFormat{
			Compressed: false,
			Format:     FormatUncompressed,
			Schema:     readDecompressedSchema(data),
		}
	}

	if len(data) >= 2 && data[0] == ZlibHeader0 && data[1] == ZlibHeader1 {
		return DetectedFormat{
			Compressed: true,
			Format:     FormatDynastyCommon,
			Schema:     readCompressedSchema(data),
		}
	}

	return DetectedFormat{
		Compressed: true,
		Format:     FormatDynasty,
		Schema:     readCompressedSchema(data),
	}
}

// readCompressedSchema reads schema major/minor from Madden-aligned offsets.
// CFB offsets may differ; update once confirmed during RE.
func readCompressedSchema(data []byte) *SchemaVersion {
	if len(data) < 0x46 {
		return nil
	}
	return &SchemaVersion{
		Major: int(binary.LittleEndian.Uint32(data[0x3e:])),
		Minor: int(binary.LittleEndian.Uint32(data[0x42:])),
	}
}

// readDecompressedSchema reads schema major/minor from Madden-aligned offsets.
func readDecompressedSchema(data []byte) *SchemaVersion {
	if len(data) < 0x30 {
		return nil
	}
	return &SchemaVersion{
		Major: int(binary.BigEndian.Uint32(data[0x2c:])),
		Minor: int(binary.BigEndian.Uint32(data[0x28:])),
	}
}
