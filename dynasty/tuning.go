package dynasty

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type skillGroupAttributeDef struct {
	name          string
	playerAbility string
}

type skillGroupBucketInfo struct {
	labels     []string
	slots      []int // upgrade slots per bucket from tuning Primary/Secondary/TertiarySkills
	attributes [][]skillGroupAttributeDef
}

type skillGroupIndex struct {
	byKey map[string]skillGroupBucketInfo // playerType|archetypeLabel -> bucket metadata
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
	idx := skillGroupIndex{byKey: make(map[string]skillGroupBucketInfo)}
	sg, ok := tf.PrimaryTableByName("PlayerSkillGroup")
	if !ok {
		return idx
	}
	if bucketTable, ok := tf.PrimaryTableByName("PlayerSkillGroupBucket"); ok {
		_ = bucketTable.ReadRecords()
	}
	if skillTable, ok := tf.PrimaryTableByName("PlayerSkill"); ok {
		_ = skillTable.ReadRecords()
	}
	if err := sg.ReadRecords(); err != nil {
		return idx
	}
	fallback := make(map[string]skillGroupBucketInfo)
	for _, record := range sg.Records {
		playerType := normalizeEnum(stringField(record, "Archetype"))
		if playerType == "" || playerType == "First" {
			continue
		}
		buckets := skillGroupBucketInfoFromRecord(tf, record)
		if len(buckets.labels) != skillGroupCapCount {
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

func skillGroupBucketInfoFromRecord(f *File, record Record) skillGroupBucketInfo {
	ref, ok := record.Get("PlayerSkillsBucket")
	if !ok || ref.Reference == nil {
		return skillGroupBucketInfo{}
	}
	var info skillGroupBucketInfo
	for _, memberRef := range f.arrayStoreMemberRefs(ref.Reference) {
		row, ok := f.RecordByReference("PlayerSkillGroupBucket", memberRef)
		if !ok {
			continue
		}
		name := stringField(row, "Name")
		if name == "" {
			continue
		}
		info.labels = append(info.labels, name)
		skills := bucketSkillsFromRow(f, row)
		info.slots = append(info.slots, len(skills))
		info.attributes = append(info.attributes, skills)
		if len(info.labels) >= skillGroupCapCount {
			break
		}
	}
	return info
}

func bucketSkillsFromRow(f *File, row Record) []skillGroupAttributeDef {
	var skills []skillGroupAttributeDef
	for _, field := range []string{"PrimarySkills", "SecondarySkills", "TertiarySkills"} {
		ref, ok := row.Get(field)
		if !ok || ref.Reference == nil {
			continue
		}
		for _, memberRef := range f.arrayStoreMemberRefs(ref.Reference) {
			skillRow, ok := f.RecordByReference("PlayerSkill", memberRef)
			if !ok {
				continue
			}
			ability := stringField(skillRow, "PlayerAbility")
			if ability == "" {
				continue
			}
			skills = append(skills, skillGroupAttributeDef{
				name:          stringField(skillRow, "Name"),
				playerAbility: ability,
			})
		}
	}
	return skills
}

func bucketSkillUpgradeSlots(f *File, row Record) int {
	return len(bucketSkillsFromRow(f, row))
}

func (idx skillGroupIndex) bucketInfo(record Record) skillGroupBucketInfo {
	if len(idx.byKey) == 0 {
		return skillGroupBucketInfo{}
	}
	playerType := normalizeEnum(stringField(record, "PlayerType"))
	label := archetypeLabelFromRecord(record)
	if buckets, ok := idx.byKey[skillGroupKey(playerType, label)]; ok {
		return buckets
	}
	if buckets, ok := idx.byKey[skillGroupKey(playerType, "")]; ok {
		return buckets
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
		return idx.byKey[keys[len(keys)-1]]
	}
	return skillGroupBucketInfo{}
}

func applySkillGroupLabels(player *PlayerExport, record Record, idx skillGroupIndex) {
	if player == nil || len(player.SkillGroupCaps) == 0 {
		return
	}
	info := idx.bucketInfo(record)
	if len(info.labels) == 0 {
		return
	}
	player.SkillGroupLabels = append([]string(nil), info.labels...)
	player.SkillGroupAttributeCounts = append([]int(nil), info.slots...)
	groups := make([]SkillGroupExport, 0, len(player.SkillGroupCaps))
	for i := range player.SkillGroupCaps {
		entry := SkillGroupExport{Slot: i + 1}
		if i < len(player.SkillGroupUnlockedSlots) {
			entry.UnlockedSlots = player.SkillGroupUnlockedSlots[i]
		}
		if i < len(player.SkillGroupCaps) {
			entry.CappedSlots = player.SkillGroupCaps[i]
		}
		if i < len(info.labels) {
			entry.Label = info.labels[i]
		}
		if i < len(info.slots) {
			entry.AttributeCount = info.slots[i]
		}
		if i < len(info.attributes) {
			for _, attr := range info.attributes[i] {
				ae := SkillGroupAttributeExport{
					Name:          attr.name,
					PlayerAbility: attr.playerAbility,
					RatingKey:     playerAbilityRatingKey(attr.playerAbility),
				}
				if player.Ratings != nil {
					if val, ok := player.Ratings[ae.RatingKey]; ok {
						v := val
						ae.Rating = &v
					}
				}
				entry.Attributes = append(entry.Attributes, ae)
			}
		}
		groups = append(groups, entry)
	}
	player.SkillGroups = groups
}

func playerAbilityRatingKey(ability string) string {
	return ability + "Rating"
}
