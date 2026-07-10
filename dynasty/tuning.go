package dynasty

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type skillGroupIndex struct {
	byKey map[string][]string // playerType|archetypeLabel -> bucket names
}

var (
	skillGroupIndexMu   sync.Mutex
	skillGroupIndexByPath = make(map[string]skillGroupIndex)
)

func (f *File) skillGroupIndex() skillGroupIndex {
	if f == nil {
		return skillGroupIndex{}
	}
	path := f.settings.TuningPath
	if path == "" {
		path = discoverTuningPath(f.settings.SchemaDir)
	}
	if path == "" {
		return skillGroupIndex{}
	}

	skillGroupIndexMu.Lock()
	if idx, ok := skillGroupIndexByPath[path]; ok {
		skillGroupIndexMu.Unlock()
		return idx
	}
	skillGroupIndexMu.Unlock()

	idx, err := loadSkillGroupIndex(path, f.settings.SchemaDir)
	if err != nil {
		return skillGroupIndex{}
	}

	skillGroupIndexMu.Lock()
	skillGroupIndexByPath[path] = idx
	skillGroupIndexMu.Unlock()
	return idx
}

func loadSkillGroupIndex(tuningPath, schemaDir string) (skillGroupIndex, error) {
	settings := DefaultSettings()
	settings.SchemaDir = schemaDir
	settings.TuningPath = tuningPath
	settings.AutoParse = true

	tf, err := Open(tuningPath, &settings)
	if err != nil {
		return skillGroupIndex{}, fmt.Errorf("cfb-dynasty: open tuning: %w", err)
	}
	idx := buildSkillGroupIndex(tf)
	// Release decoded tuning tables; only the small label index is retained.
	for i := range tf.tables {
		tf.tables[i].Records = nil
	}
	return idx, nil
}

func discoverTuningPath(schemaDir string) string {
	if schemaDir == "" {
		return ""
	}
	candidates := []string{
		filepath.Join(schemaDir, "cfb27-db-data", "2", "dynasty-tuning-binary.FTC"),
		filepath.Join(schemaDir, "cfb27-db-data", "0", "dynasty-tuning-binary.FTC"),
		filepath.Join(schemaDir, "cfb27-db-data", "1", "dynasty-tuning-binary.FTC"),
		filepath.Join(schemaDir, "dynasty-tuning-binary.FTC"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func buildSkillGroupIndex(tf *File) skillGroupIndex {
	idx := skillGroupIndex{byKey: make(map[string][]string)}
	sg, ok := tf.PrimaryTableByName("PlayerSkillGroup")
	if !ok {
		return idx
	}
	if bucketTable, ok := tf.PrimaryTableByName("PlayerSkillGroupBucket"); ok {
		_ = bucketTable.ReadRecords()
	}
	if err := sg.ReadRecords(); err != nil {
		return idx
	}
	fallback := make(map[string][]string)
	for _, record := range sg.Records {
		playerType := normalizeEnum(stringField(record, "Archetype"))
		if playerType == "" || playerType == "First" {
			continue
		}
		buckets := skillGroupBucketsFromRecord(tf, record)
		if len(buckets) != skillGroupCapCount {
			continue
		}
		label := stringField(record, "Name")
		if label != "" {
			idx.byKey[skillGroupKey(playerType, label)] = buckets
		} else if _, exists := fallback[playerType]; !exists {
			fallback[playerType] = buckets
		}
	}
	for playerType, buckets := range fallback {
		key := skillGroupKey(playerType, "")
		if _, exists := idx.byKey[key]; !exists {
			idx.byKey[key] = buckets
		}
	}
	return idx
}

func skillGroupKey(playerType, label string) string {
	return playerType + "|" + label
}

func skillGroupBucketsFromRecord(f *File, record Record) []string {
	ref, ok := record.Get("PlayerSkillsBucket")
	if !ok || ref.Reference == nil {
		return nil
	}
	var buckets []string
	for _, memberRef := range f.arrayStoreMemberRefs(ref.Reference) {
		row, ok := f.RecordByReference("PlayerSkillGroupBucket", memberRef)
		if !ok {
			continue
		}
		name := stringField(row, "Name")
		if name == "" {
			continue
		}
		buckets = append(buckets, name)
	}
	return buckets
}

func (idx skillGroupIndex) bucketLabels(record Record) []string {
	if len(idx.byKey) == 0 {
		return nil
	}
	playerType := normalizeEnum(stringField(record, "PlayerType"))
	label := archetypeLabelFromRecord(record)
	if buckets, ok := idx.byKey[skillGroupKey(playerType, label)]; ok {
		return append([]string(nil), buckets...)
	}
	if buckets, ok := idx.byKey[skillGroupKey(playerType, "")]; ok {
		return append([]string(nil), buckets...)
	}
	prefix := playerType + "|"
	var keys []string
	for key := range idx.byKey {
		if strings.HasPrefix(key, prefix) {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	if len(keys) > 0 {
		return append([]string(nil), idx.byKey[keys[len(keys)-1]]...)
	}
	return nil
}

func applySkillGroupLabels(player *PlayerExport, record Record, idx skillGroupIndex) {
	if player == nil || len(player.SkillGroupCaps) == 0 {
		return
	}
	labels := idx.bucketLabels(record)
	if len(labels) == 0 {
		return
	}
	player.SkillGroupLabels = labels
	groups := make([]SkillGroupExport, 0, len(player.SkillGroupCaps))
	for i, cap := range player.SkillGroupCaps {
		entry := SkillGroupExport{Cap: cap}
		if i < len(labels) {
			entry.Label = labels[i]
		}
		groups = append(groups, entry)
	}
	player.SkillGroups = groups
}
