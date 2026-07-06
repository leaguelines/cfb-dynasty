package dynasty

// Known byte signatures observed in Madden franchise saves and suspected in CFB.
const (
	MagicFrTk = "FrTk" // 46 72 54 6b — decompressed FranTk header
	MagicSPBF = "SPBF" // table marker
	MagicASTO = "ASTO" // alternate table marker
	MagicSPEX = "SPEX" // alternate table marker

	ZlibHeader0 = 0x78
	ZlibHeader1 = 0x9c

	// CompressedDataOffset is Madden's zlib payload skip; CFB may differ.
	CompressedDataOffset = 0x52

	// TableMarkerOffset is the byte offset of the SPBF/ASTO/SPEX marker in a table header.
	TableMarkerOffset = 0x94

	// table2StringIndexSentinel marks an indirect table2 string reference; the real
	// index is stored in another field (CFB Team LongName uses DIV_SLOTNUMBER).
	table2StringIndexSentinel = 0x80000000

	DefaultGameYear = 27
)

// Format describes the on-disk container variant.
type Format string

const (
	FormatUnknown         Format = "unknown"
	FormatDynastyCommon   Format = "dynasty-common" // suspected small/common compressed saves
	FormatDynasty         Format = "dynasty"        // suspected full dynasty saves
	FormatUncompressed    Format = "uncompressed"     // FrTk header present
)

// TableMarker identifies a table header signature in unpacked bytes.
type TableMarker string

const (
	TableMarkerSPBF TableMarker = MagicSPBF
	TableMarkerASTO TableMarker = MagicASTO
	TableMarkerSPEX TableMarker = MagicSPEX
)

// AllTableMarkers lists markers scanned during table discovery.
var AllTableMarkers = []TableMarker{
	TableMarkerSPBF,
	TableMarkerASTO,
	TableMarkerSPEX,
}
