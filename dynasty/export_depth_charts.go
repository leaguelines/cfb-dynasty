package dynasty

import "sort"

// buildDepthChartExports groups forced depth chart entries by team.
func (f *File) buildDepthChartExports() ([]DepthChartExport, error) {
	depthTable, ok := f.PrimaryTableByName("ForcedDepthChartEntry")
	if !ok {
		return nil, nil
	}
	if err := depthTable.ReadRecords(); err != nil {
		return nil, err
	}

	teams := f.teamMaps()
	byTeam := make(map[int]*DepthChartExport)

	for _, record := range depthTable.Records {
		player, playerID, ok := f.playerRecordFromField(record, "Player")
		if !ok {
			continue
		}
		row, ok := intFieldOK(player, "TeamIndex")
		if !ok {
			continue
		}
		teamID, ok := teams.exportID(row)
		if !ok {
			continue
		}

		chart, ok := byTeam[teamID]
		if !ok {
			chart = &DepthChartExport{
				TeamID:   teamID,
				TeamName: teams.nameFromID(teamID),
			}
			byTeam[teamID] = chart
		}

		slot := DepthChartSlotExport{
			Position:     stringField(record, "Position"),
			Depth:        intField(record, "CurrentDepth"),
			PlayerID:     playerID,
			FirstName:    stringField(player, "FirstName"),
			LastName:     stringField(player, "LastName"),
			UserEditable: boolField(record, "UserEditable"),
		}
		if locked, ok := intFieldOK(record, "LockedDepth"); ok && locked > 0 {
			slot.LockedDepth = &locked
		}
		chart.Slots = append(chart.Slots, slot)
	}

	exports := make([]DepthChartExport, 0, len(byTeam))
	for _, chart := range byTeam {
		sortDepthChartSlots(chart.Slots)
		exports = append(exports, *chart)
	}
	sort.Slice(exports, func(i, j int) bool {
		return exports[i].TeamID < exports[j].TeamID
	})
	return exports, nil
}

func sortDepthChartSlots(slots []DepthChartSlotExport) {
	sort.Slice(slots, func(i, j int) bool {
		if slots[i].Position != slots[j].Position {
			return slots[i].Position < slots[j].Position
		}
		if slots[i].Depth != slots[j].Depth {
			return slots[i].Depth < slots[j].Depth
		}
		return slots[i].LastName < slots[j].LastName
	})
}
