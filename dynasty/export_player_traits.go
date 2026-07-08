package dynasty

import (
	"sort"
	"strings"
)

// archetypeTraitsFromRecord returns labels for every active PT_* trait flag on the player.
func archetypeTraitsFromRecord(record Record) []string {
	labels := make(map[string]struct{})
	for _, entry := range archetypePTEntries {
		if boolField(record, entry.field) {
			label := entry.label
			if entry.field == "PT_DLRUNSTOPPER" {
				label = dlRunStopperLabel(record)
			}
			labels[label] = struct{}{}
		}
	}

	for key, value := range record.Fields {
		if !strings.HasPrefix(key, "PT_") || !value.Bool {
			continue
		}
		if label := archetypeTraitLabel(key); label != "" {
			labels[label] = struct{}{}
		}
	}

	if len(labels) == 0 {
		return nil
	}
	out := make([]string, 0, len(labels))
	for label := range labels {
		out = append(out, label)
	}
	sort.Strings(out)
	return out
}

func archetypeTraitLabel(field string) string {
	for _, entry := range archetypePTEntries {
		if entry.field == field {
			return entry.label
		}
	}
	name := strings.TrimPrefix(field, "PT_")
	name = strings.TrimSuffix(name, "_")
	if name == "" || strings.HasPrefix(name, "UNUSED") {
		return ""
	}
	return name
}
