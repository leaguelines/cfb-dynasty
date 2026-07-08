package dynasty

import "sort"

func (f *File) buildBowlGameExports() ([]BowlGameExport, error) {
	table, ok := f.PrimaryTableByName("BowlGame")
	if !ok {
		return nil, nil
	}
	if err := table.ReadRecords(); err != nil {
		return nil, err
	}
	exports := make([]BowlGameExport, 0, table.ActiveRecordCount())
	for _, record := range table.Records {
		name := stringField(record, "Name")
		if name == "" {
			continue
		}
		export := BowlGameExport{
			ID:            record.Index,
			Name:          name,
			IsPlayoffBowl: boolField(record, "IsPlayoffBowl"),
		}
		setOptionalPositiveInt(record, "PlayoffBracketSlot", &export.PlayoffBracketSlot)
		exports = append(exports, export)
	}
	sort.Slice(exports, func(i, j int) bool { return exports[i].ID < exports[j].ID })
	return exports, nil
}
