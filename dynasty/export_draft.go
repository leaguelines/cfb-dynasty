package dynasty

import "sort"

func (f *File) buildDraftPickExports() ([]DraftPickExport, error) {
	table, ok := f.PrimaryTableByName("DraftPick")
	if !ok {
		return nil, nil
	}
	if err := table.ReadRecords(); err != nil {
		return nil, err
	}
	teams := f.teamMaps()
	exports := make([]DraftPickExport, 0, table.ActiveRecordCount())
	for _, record := range table.Records {
		export := DraftPickExport{
			ID:            record.Index,
			PositionGroup: normalizeEnum(stringField(record, "PositionGroup")),
		}
		setOptionalPositiveInt(record, "YearOffset", &export.YearOffset)
		setOptionalPositiveInt(record, "Round", &export.Round)
		setOptionalPositiveInt(record, "PickNumber", &export.PickNumber)
		export.TeamID, export.TeamName = teamRefExport(f, record, "Team", teams)
		exports = append(exports, export)
	}
	sort.Slice(exports, func(i, j int) bool { return exports[i].ID < exports[j].ID })
	return exports, nil
}
