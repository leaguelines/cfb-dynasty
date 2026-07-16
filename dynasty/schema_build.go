package dynasty

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	_ "embed"
)

//go:embed data/extra-schemas.json
var embeddedExtraSchemasJSON []byte

// DefaultCFB27SchemaMajor is the FranTk schema major used by current CFB 27
// gzip bundles (for example C27_468_2.gz). Split Frosty extracts do not embed
// dataMajorVersion, so builders fall back to this when game year is 27.
const DefaultCFB27SchemaMajor = 468

// SchemaBuildOptions controls FTX → gzip schema bundle generation.
type SchemaBuildOptions struct {
	// Source is a patch folder (…/cfb27-db-data/2) or the parent cfb27-db-data
	// directory (newest revision is chosen).
	Source string

	// OutDir is where C{year}_{major}_{minor}.gz is written.
	OutDir string

	// Major/Minor/GameYear override auto-detection when non-nil / non-zero.
	Major    *int
	Minor    *int
	GameYear *int

	// IncludeExtras merges madden-franchise core extras (User, Stage, …).
	// Defaults to true when unset via BuildSchemaBundle.
	IncludeExtras *bool
}

// SchemaBuildResult is the written bundle metadata.
type SchemaBuildResult struct {
	Path     string
	Major    int
	Minor    int
	GameYear int
	Tables   int
	Enums    int
	// MinorSource describes how minor was chosen (revision, directory, flag).
	MinorSource string
	// MajorSource describes how major was chosen.
	MajorSource string
}

type ftxRoot struct {
	XMLName             xml.Name   `xml:"FranTkData"`
	FileName            string     `xml:"fileName,attr"`
	Namespace           string     `xml:"namespace,attr"`
	DatabaseName        string     `xml:"databaseName,attr"`
	DataMajorVersion    string     `xml:"dataMajorVersion,attr"`
	DataMinorVersion    string     `xml:"dataMinorVersion,attr"`
	DataRevisionVersion string     `xml:"dataRevisionVersion,attr"`
	Schemas             ftxSchemas `xml:"schemas"`
}

type ftxSchemas struct {
	Schema []ftxSchema `xml:"schema"`
	Enum   []ftxEnum   `xml:"enum"`
}

type ftxSchema struct {
	Name       string         `xml:"name,attr"`
	NumMembers string         `xml:"numMembers,attr"`
	Base       string         `xml:"base,attr"`
	AssetID    string         `xml:"assetId,attr"`
	Attributes []ftxAttribute `xml:"attribute"`
}

type ftxEnum struct {
	Name       string         `xml:"name,attr"`
	AssetID    string         `xml:"assetId,attr"`
	Attributes []ftxAttribute `xml:"attribute"`
}

type ftxAttribute struct {
	Name      string `xml:"name,attr"`
	Idx       string `xml:"idx,attr"`
	Type      string `xml:"type,attr"`
	MinValue  string `xml:"minValue,attr"`
	MaxValue  string `xml:"maxValue,attr"`
	MaxLen    string `xml:"maxLen,attr"`
	Default   string `xml:"default,attr"`
	Final     string `xml:"final,attr"`
	Const     string `xml:"const,attr"`
	Value     string `xml:"value,attr"`
}

type buildEnum struct {
	Name               string
	AssetID            string
	IsRecordPersistent string
	Members            []buildEnumMember
	MaxLength          int
}

type buildEnumMember struct {
	Name             string
	Index            int
	Value            int
	UnformattedValue string
}

type buildTable struct {
	AssetID    string
	Name       string
	Base       string
	NumMembers string
	Attributes []buildAttr
}

type buildAttr struct {
	Index     string
	Name      string
	Type      string
	MinValue  string
	MaxValue  string
	MaxLength string
	Default   string
	Final     string
	Const     string
	Enum      *buildEnum
}

var (
	cfbYearInPath = regexp.MustCompile(`(?i)cfb(\d{2})`)
	numericDir    = regexp.MustCompile(`^\d+$`)
)

// BuildSchemaBundle evaluates FranTk .FTX files under opts.Source and writes a
// gzip JSON schema bundle compatible with LoadSchemaFile.
func BuildSchemaBundle(opts SchemaBuildOptions) (*SchemaBuildResult, error) {
	if strings.TrimSpace(opts.Source) == "" {
		return nil, fmt.Errorf("%w: schema-build source required", ErrInvalidSchema)
	}
	source, err := resolveSchemaSource(opts.Source)
	if err != nil {
		return nil, err
	}

	meta, majorSource, minorSource, err := detectSchemaMeta(source)
	if err != nil {
		return nil, err
	}
	if opts.Major != nil {
		meta.Major = *opts.Major
		majorSource = "flag"
	} else if meta.Major == 0 {
		year := meta.GameYear
		if opts.GameYear != nil {
			year = *opts.GameYear
		}
		if year == 0 {
			year = detectGameYear(source)
		}
		if year == 27 {
			meta.Major = DefaultCFB27SchemaMajor
			majorSource = "default-cfb27"
		} else {
			return nil, fmt.Errorf("%w: cannot detect schema major in %q; pass -major", ErrInvalidSchema, source)
		}
	}
	if opts.Minor != nil {
		meta.Minor = *opts.Minor
		minorSource = "flag"
	} else if meta.Minor < 0 {
		return nil, fmt.Errorf("%w: cannot detect schema minor in %q; pass -minor", ErrInvalidSchema, source)
	}
	if opts.GameYear != nil {
		meta.GameYear = *opts.GameYear
	} else if meta.GameYear == 0 {
		meta.GameYear = detectGameYear(source)
	}
	if meta.GameYear == 0 {
		meta.GameYear = 27
	}

	includeExtras := true
	if opts.IncludeExtras != nil {
		includeExtras = *opts.IncludeExtras
	}

	tables, enums, err := evaluateFTXDir(source, includeExtras)
	if err != nil {
		return nil, err
	}

	outDir := opts.OutDir
	if outDir == "" {
		outDir = "."
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, fmt.Errorf("cfb-dynasty: create output dir %q: %w", outDir, err)
	}
	name := fmt.Sprintf("C%d_%d_%d.gz", meta.GameYear, meta.Major, meta.Minor)
	outPath := filepath.Join(outDir, name)
	if err := writeSchemaGzip(outPath, meta, tables); err != nil {
		return nil, err
	}

	return &SchemaBuildResult{
		Path:        outPath,
		Major:       meta.Major,
		Minor:       meta.Minor,
		GameYear:    meta.GameYear,
		Tables:      len(tables),
		Enums:       len(enums),
		MinorSource: minorSource,
		MajorSource: majorSource,
	}, nil
}

type detectedMeta struct {
	Major    int
	Minor    int
	GameYear int
}

func resolveSchemaSource(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("cfb-dynasty: schema source %q: %w", path, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%w: schema source must be a directory", ErrInvalidSchema)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// Parent cfb27-db-data: prefer numbered patch children when present.
	entries, err := os.ReadDir(abs)
	if err != nil {
		return "", err
	}
	type cand struct {
		path     string
		revision int
		dirNum   int
	}
	var cands []cand
	for _, e := range entries {
		if !e.IsDir() || !numericDir.MatchString(e.Name()) {
			continue
		}
		child := filepath.Join(abs, e.Name())
		if !hasFTXTree(child) {
			continue
		}
		rev, _, _ := detectRevision(child)
		dirNum, _ := strconv.Atoi(e.Name())
		cands = append(cands, cand{path: child, revision: rev, dirNum: dirNum})
	}
	if len(cands) > 0 {
		sort.Slice(cands, func(i, j int) bool {
			if cands[i].revision != cands[j].revision {
				return cands[i].revision > cands[j].revision
			}
			return cands[i].dirNum > cands[j].dirNum
		})
		return cands[0].path, nil
	}

	if hasFTXTree(abs) {
		return abs, nil
	}
	return "", fmt.Errorf("%w: no FTX schema tree under %q", ErrInvalidSchema, path)
}

func hasFTXTree(dir string) bool {
	found := false
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Ext(d.Name()), ".ftx") {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found
}

func detectSchemaMeta(dir string) (meta detectedMeta, majorSource, minorSource string, err error) {
	meta.Minor = -1

	// Prefer FranTk.College / franchise-schemas.FTX root meta (not Core or Football).
	bestScore := -1 << 30
	var best *ftxRoot
	var bestPath string
	walkErr := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.EqualFold(filepath.Ext(d.Name()), ".ftx") {
			return nil
		}
		root, perr := parseFTXFile(path)
		if perr != nil || root == nil {
			return nil
		}
		if root.DataMajorVersion == "" && root.DataMinorVersion == "" && root.DatabaseName == "" {
			return nil
		}
		score := scoreFTXVersionSource(path, root)
		if score > bestScore {
			bestScore = score
			best = root
			bestPath = path
		}
		return nil
	})
	if walkErr != nil {
		return meta, "", "", walkErr
	}

	if best != nil && bestScore > 0 {
		if best.DataMajorVersion != "" {
			if n, e := strconv.Atoi(best.DataMajorVersion); e == nil {
				meta.Major = n
				majorSource = "franchise-dataMajorVersion"
				if !isCollegeFranchiseRoot(bestPath, best) {
					majorSource = "ftx-dataMajorVersion"
				}
			}
		}
		if best.DataMinorVersion != "" {
			if n, e := strconv.Atoi(best.DataMinorVersion); e == nil {
				meta.Minor = n
				minorSource = "franchise-dataMinorVersion"
				if !isCollegeFranchiseRoot(bestPath, best) {
					minorSource = "ftx-dataMinorVersion"
				}
			}
		}
		if best.DatabaseName != "" {
			meta.GameYear = gameYearFromDatabaseName(best.DatabaseName)
		}
	}

	// Older split extracts omit root major/minor — fall back to revision / directory.
	if meta.Minor < 0 {
		rev, revCount, rerr := detectRevision(dir)
		if rerr != nil {
			return meta, "", "", rerr
		}
		if revCount > 0 {
			meta.Minor = rev
			minorSource = "dataRevisionVersion"
		}
	}
	if meta.Minor < 0 {
		base := filepath.Base(dir)
		if numericDir.MatchString(base) {
			n, _ := strconv.Atoi(base)
			meta.Minor = n
			minorSource = "directory"
		}
	}
	if meta.GameYear == 0 {
		meta.GameYear = detectGameYear(dir)
	}
	return meta, majorSource, minorSource, nil
}

// scoreFTXVersionSource ranks FranTkData roots for bundle meta. College/franchise
// roots win over Football and Core (which ship their own major/minor).
func scoreFTXVersionSource(path string, root *ftxRoot) int {
	base := strings.ToLower(filepath.Base(path))
	ns := strings.ToLower(root.Namespace)
	fileName := strings.ToLower(root.FileName)
	score := 0

	switch {
	case base == "franchise-schemas.ftx":
		score += 200
	case ns == "frantk.college":
		score += 150
	case fileName == "franchise-schemas" || strings.HasPrefix(fileName, "franchise-schemas"):
		// Root aggregator only — nested paths look like Franchise-Schemas\Player.ftx.
		if !strings.Contains(fileName, `\`) && !strings.Contains(fileName, "/") {
			score += 120
		}
	}

	switch {
	case base == "core-schemas.ftx", ns == "frantk.core", fileName == "core-schemas":
		score -= 200
	case base == "football-schemas.ftx", ns == "frantk.football", fileName == "football-schemas":
		score -= 100
	}

	if root.DataMajorVersion != "" {
		score += 10
	}
	if root.DatabaseName != "" {
		score += 5
	}
	return score
}

func isCollegeFranchiseRoot(path string, root *ftxRoot) bool {
	base := strings.ToLower(filepath.Base(path))
	ns := strings.ToLower(root.Namespace)
	fileName := strings.ToLower(root.FileName)
	if base == "franchise-schemas.ftx" || ns == "frantk.college" {
		return true
	}
	if (fileName == "franchise-schemas" || strings.HasPrefix(fileName, "franchise-schemas")) &&
		!strings.Contains(fileName, `\`) && !strings.Contains(fileName, "/") {
		return true
	}
	return false
}

func detectRevision(dir string) (revision, count int, err error) {
	counts := map[int]int{}
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.EqualFold(filepath.Ext(d.Name()), ".ftx") {
			return nil
		}
		// Only need the opening tag.
		data, rerr := os.ReadFile(path)
		if rerr != nil {
			return nil
		}
		head := data
		if len(head) > 2048 {
			head = head[:2048]
		}
		m := regexp.MustCompile(`dataRevisionVersion="(\d+)"`).FindSubmatch(head)
		if m == nil {
			return nil
		}
		n, _ := strconv.Atoi(string(m[1]))
		counts[n]++
		return nil
	})
	if err != nil {
		return 0, 0, err
	}
	best, bestN := -1, 0
	for rev, n := range counts {
		if n > bestN || (n == bestN && rev > best) {
			best, bestN = rev, n
		}
	}
	return best, bestN, nil
}

func detectGameYear(path string) int {
	if m := cfbYearInPath.FindStringSubmatch(path); m != nil {
		n, _ := strconv.Atoi(m[1])
		return n
	}
	return 0
}

func gameYearFromDatabaseName(name string) int {
	for _, prefix := range []string{"CollegeFB", "Madden"} {
		re := regexp.MustCompile(prefix + `(\d{2})`)
		if m := re.FindStringSubmatch(name); m != nil {
			n, _ := strconv.Atoi(m[1])
			return n
		}
	}
	return 0
}

func evaluateFTXDir(dir string, includeExtras bool) ([]buildTable, map[string]*buildEnum, error) {
	var paths []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.EqualFold(filepath.Ext(d.Name()), ".ftx") {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("cfb-dynasty: walk FTX %q: %w", dir, err)
	}
	sort.Strings(paths)

	enumMap := map[string]*buildEnum{}
	var enums []*buildEnum
	tableMap := map[string]*buildTable{}
	var tables []*buildTable

	for _, path := range paths {
		root, err := parseFTXFile(path)
		if err != nil {
			return nil, nil, err
		}
		if root == nil {
			continue
		}
		for _, e := range root.Schemas.Enum {
			en := parseBuildEnum(e)
			if en.Name == "" {
				continue
			}
			enumMap[en.Name] = en
			enums = append(enums, en)
		}
		for _, s := range root.Schemas.Schema {
			t := parseBuildTable(s)
			if t.Name == "" {
				continue
			}
			tableMap[t.Name] = t
			tables = append(tables, t)
		}
	}

	for _, en := range enums {
		setEnumMemberLength(en)
	}

	if includeExtras {
		extras, err := loadExtraSchemas()
		if err != nil {
			return nil, nil, err
		}
		for _, t := range extras {
			if _, exists := tableMap[t.Name]; exists {
				continue
			}
			// Resolve string enum refs after FTX enums are known.
			for i := range t.Attributes {
				attr := &t.Attributes[i]
				if attr.Enum == nil && attr.Type != "" {
					if en, ok := enumMap[attr.Type]; ok {
						attr.Enum = en
					}
				}
			}
			tableMap[t.Name] = t
			tables = append([]*buildTable{t}, tables...)
		}
	}

	for _, t := range tables {
		for i := range t.Attributes {
			attr := &t.Attributes[i]
			if attr.Enum != nil {
				continue
			}
			if en, ok := enumMap[attr.Type]; ok {
				attr.Enum = en
			}
		}
	}

	calculateInheritedSchemas(tables)
	out := make([]buildTable, 0, len(tables))
	for _, t := range tables {
		out = append(out, *t)
	}
	return out, enumMap, nil
}

func parseFTXFile(path string) (*ftxRoot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cfb-dynasty: read FTX %q: %w", path, err)
	}
	// Skip non-schema legacy files that aren't FranTkData.
	if !bytes.Contains(data, []byte("<FranTkData")) {
		return nil, nil
	}
	var root ftxRoot
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("cfb-dynasty: parse FTX %q: %w", path, err)
	}
	return &root, nil
}

func parseBuildEnum(e ftxEnum) *buildEnum {
	en := &buildEnum{Name: e.Name, AssetID: e.AssetID}
	for _, a := range e.Attributes {
		idx, _ := strconv.Atoi(a.Idx)
		val, _ := strconv.Atoi(a.Value)
		en.Members = append(en.Members, buildEnumMember{
			Name:             a.Name,
			Index:            idx,
			Value:            val,
			UnformattedValue: enumValueToBinary(val),
		})
	}
	return en
}

func parseBuildTable(s ftxSchema) *buildTable {
	t := &buildTable{
		AssetID:    s.AssetID,
		Name:       s.Name,
		Base:       s.Base,
		NumMembers: s.NumMembers,
	}
	for i, a := range s.Attributes {
		idx := a.Idx
		if idx == "" {
			idx = strconv.Itoa(i)
		}
		t.Attributes = append(t.Attributes, buildAttr{
			Index:     idx,
			Name:      a.Name,
			Type:      a.Type,
			MinValue:  a.MinValue,
			MaxValue:  a.MaxValue,
			MaxLength: a.MaxLen,
			Default:   decodeXMLDefault(a.Default),
			Final:     a.Final,
			Const:     a.Const,
		})
	}
	return t
}

func decodeXMLDefault(v string) string {
	if v == "" {
		return ""
	}
	replacer := strings.NewReplacer(
		"&#xD;", "\r",
		"&#xA;", "\n",
		"&amp;", "&",
		"&gt;", ">",
		"&lt;", "<",
		"&quot;", `"`,
	)
	return replacer.Replace(v)
}

func enumValueToBinary(value int) string {
	if value < 0 {
		return strconv.FormatUint(uint64(uint32(-(value) - 1)), 2)
	}
	return strconv.FormatUint(uint64(uint32(value)), 2)
}

func setEnumMemberLength(en *buildEnum) {
	if en == nil || len(en.Members) == 0 {
		return
	}
	maxVal := en.Members[0]
	hasNeg := false
	for _, m := range en.Members {
		if m.Value < 0 {
			hasNeg = true
		}
		if m.Value > maxVal.Value {
			maxVal = m
		}
	}
	maxLen := len(maxVal.UnformattedValue)
	if hasNeg {
		maxLen++
	}
	en.MaxLength = maxLen
	for i := range en.Members {
		m := &en.Members[i]
		if m.Value < 0 {
			// Match madden-franchise: '1' + absBits padded to maxLen-1.
			m.UnformattedValue = "1" + padLeft(m.UnformattedValue, maxLen-1, '0')
		} else {
			m.UnformattedValue = padLeft(m.UnformattedValue, maxLen, '0')
		}
	}
}

func padLeft(s string, n int, pad byte) string {
	if len(s) >= n {
		return s
	}
	return strings.Repeat(string(pad), n-len(s)) + s
}

func calculateInheritedSchemas(tables []*buildTable) {
	byName := make(map[string]*buildTable, len(tables))
	for _, t := range tables {
		byName[t.Name] = t
	}
	for _, schema := range tables {
		if schema.Base == "" || strings.Contains(schema.Base, "()") {
			continue
		}
		base := byName[schema.Base]
		if base == nil {
			continue
		}
		for index, baseAttr := range base.Attributes {
			oldIndex := -1
			for i, a := range schema.Attributes {
				if a.Name == baseAttr.Name {
					oldIndex = i
					break
				}
			}
			if oldIndex < 0 {
				continue
			}
			arrayMoveAttrs(schema.Attributes, oldIndex, index)
		}
	}
}

func arrayMoveAttrs(attrs []buildAttr, oldIndex, newIndex int) {
	if oldIndex == newIndex || oldIndex < 0 || oldIndex >= len(attrs) {
		return
	}
	item := attrs[oldIndex]
	if oldIndex < newIndex {
		copy(attrs[oldIndex:newIndex], attrs[oldIndex+1:newIndex+1])
		attrs[newIndex] = item
		return
	}
	copy(attrs[newIndex+1:oldIndex+1], attrs[newIndex:oldIndex])
	attrs[newIndex] = item
}

func loadExtraSchemas() ([]*buildTable, error) {
	var raw []map[string]any
	if err := json.Unmarshal(embeddedExtraSchemasJSON, &raw); err != nil {
		return nil, fmt.Errorf("%w: extra-schemas.json: %v", ErrInvalidSchema, err)
	}
	var out []*buildTable
	for _, item := range raw {
		if typ, _ := item["type"].(string); typ == "enum" {
			continue
		}
		name, _ := item["name"].(string)
		if name == "" {
			continue
		}
		t := &buildTable{
			Name:       name,
			Base:       anyString(item["base"]),
			NumMembers: anyString(item["numMembers"]),
			AssetID:    anyString(item["assetId"]),
		}
		attrs, _ := item["attributes"].([]any)
		for i, a := range attrs {
			am, ok := a.(map[string]any)
			if !ok {
				continue
			}
			idx := anyString(am["index"])
			if idx == "" {
				idx = anyString(am["idx"])
			}
			if idx == "" {
				idx = strconv.Itoa(i)
			}
			attr := buildAttr{
				Index:     idx,
				Name:      anyString(am["name"]),
				Type:      anyString(am["type"]),
				MinValue:  anyString(am["minValue"]),
				MaxValue:  anyString(am["maxValue"]),
				MaxLength: anyString(am["maxLength"]),
				Default:   anyString(am["default"]),
				Final:     anyString(am["final"]),
				Const:     anyString(am["const"]),
			}
			switch e := am["enum"].(type) {
			case string:
				if e != "" {
					attr.Type = firstNonEmpty(attr.Type, e)
					// Deferred resolve by type name.
				}
			case map[string]any:
				attr.Enum = parseExtraEnumObject(e)
			}
			t.Attributes = append(t.Attributes, attr)
		}
		out = append(out, t)
	}
	return out, nil
}

func parseExtraEnumObject(m map[string]any) *buildEnum {
	en := &buildEnum{
		Name:               anyString(m["_name"]),
		AssetID:            anyString(m["_assetId"]),
		IsRecordPersistent: anyString(m["_isRecordPersistent"]),
	}
	if n, ok := m["_maxLength"].(float64); ok {
		en.MaxLength = int(n)
	}
	members, _ := m["_members"].([]any)
	for _, mem := range members {
		mm, ok := mem.(map[string]any)
		if !ok {
			continue
		}
		en.Members = append(en.Members, buildEnumMember{
			Name:             anyString(mm["_name"]),
			Index:            anyInt(mm["_index"]),
			Value:            anyInt(mm["_value"]),
			UnformattedValue: anyString(mm["_unformattedValue"]),
		})
	}
	if en.MaxLength == 0 {
		setEnumMemberLength(en)
	}
	return en
}

func writeSchemaGzip(path string, meta detectedMeta, tables []buildTable) error {
	type emitMember struct {
		Name             string `json:"_name"`
		Index            int    `json:"_index"`
		Value            int    `json:"_value"`
		UnformattedValue string `json:"_unformattedValue"`
	}
	type emitEnum struct {
		Name               string       `json:"_name"`
		AssetID            string       `json:"_assetId,omitempty"`
		IsRecordPersistent string       `json:"_isRecordPersistent,omitempty"`
		Members            []emitMember `json:"_members"`
		MaxLength          int          `json:"_maxLength"`
	}
	type emitAttr struct {
		Index     string    `json:"index,omitempty"`
		Name      string    `json:"name"`
		Type      string    `json:"type,omitempty"`
		MinValue  string    `json:"minValue,omitempty"`
		MaxValue  string    `json:"maxValue,omitempty"`
		MaxLength string    `json:"maxLength,omitempty"`
		Default   string    `json:"default,omitempty"`
		Final     string    `json:"final,omitempty"`
		Const     string    `json:"const,omitempty"`
		Enum      *emitEnum `json:"enum,omitempty"`
	}
	type emitTable struct {
		AssetID    string     `json:"assetId,omitempty"`
		Name       string     `json:"name"`
		Base       string     `json:"base,omitempty"`
		NumMembers string     `json:"numMembers,omitempty"`
		Attributes []emitAttr `json:"attributes"`
	}
	type emitBundle struct {
		Meta struct {
			Major    int `json:"major"`
			Minor    int `json:"minor"`
			GameYear int `json:"gameYear"`
		} `json:"meta"`
		Schemas   []emitTable          `json:"schemas"`
		SchemaMap map[string]emitTable `json:"schemaMap"`
	}

	toEmitEnum := func(en *buildEnum) *emitEnum {
		if en == nil {
			return nil
		}
		out := &emitEnum{
			Name:               en.Name,
			AssetID:            en.AssetID,
			IsRecordPersistent: en.IsRecordPersistent,
			MaxLength:          en.MaxLength,
		}
		for _, m := range en.Members {
			out.Members = append(out.Members, emitMember{
				Name:             m.Name,
				Index:            m.Index,
				Value:            m.Value,
				UnformattedValue: m.UnformattedValue,
			})
		}
		return out
	}

	bundle := emitBundle{SchemaMap: map[string]emitTable{}}
	bundle.Meta.Major = meta.Major
	bundle.Meta.Minor = meta.Minor
	bundle.Meta.GameYear = meta.GameYear

	// Extras appear first in tables; schemaMap only gets FTX-origin tables
	// (those present before extras were prepended). Match shipped C27 bundles:
	// schemaMap size ≈ FTX schema count, schemas includes extras.
	extraNames := map[string]struct{}{}
	if raw, err := loadExtraSchemas(); err == nil {
		for _, t := range raw {
			extraNames[t.Name] = struct{}{}
		}
	}

	for _, t := range tables {
		et := emitTable{
			AssetID:    t.AssetID,
			Name:       t.Name,
			Base:       t.Base,
			NumMembers: t.NumMembers,
		}
		for _, a := range t.Attributes {
			et.Attributes = append(et.Attributes, emitAttr{
				Index:     a.Index,
				Name:      a.Name,
				Type:      a.Type,
				MinValue:  a.MinValue,
				MaxValue:  a.MaxValue,
				MaxLength: a.MaxLength,
				Default:   a.Default,
				Final:     a.Final,
				Const:     a.Const,
				Enum:      toEmitEnum(a.Enum),
			})
		}
		bundle.Schemas = append(bundle.Schemas, et)
		if _, isExtra := extraNames[t.Name]; !isExtra {
			bundle.SchemaMap[t.Name] = et
		}
	}

	raw, err := json.Marshal(bundle)
	if err != nil {
		return fmt.Errorf("cfb-dynasty: marshal schema bundle: %w", err)
	}
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write(raw); err != nil {
		_ = zw.Close()
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("cfb-dynasty: write schema %q: %w", path, err)
	}
	return nil
}

func anyString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case float64:
		if t == float64(int(t)) {
			return strconv.Itoa(int(t))
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(t)
	default:
		return fmt.Sprint(t)
	}
}

func anyInt(v any) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case string:
		n, _ := strconv.Atoi(t)
		return n
	case int:
		return t
	default:
		return 0
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
