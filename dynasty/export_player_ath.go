package dynasty

// playerAlternatePositions maps Player row index to recruit alternate positions.
func (f *File) playerAlternatePositions() (map[int][2]string, error) {
	recruitTable, ok := f.PrimaryTableByName("Recruit")
	if !ok {
		return nil, nil
	}
	if err := recruitTable.ReadRecords(); err != nil {
		return nil, err
	}

	out := make(map[int][2]string)
	for _, record := range recruitTable.Records {
		ref, ok := record.Get("Player")
		if !ok || ref.Reference == nil {
			continue
		}
		playerRecord, _, ok := f.playerRecordByReference(ref.Reference)
		if !ok {
			continue
		}
		out[playerRecord.Index] = [2]string{
			stringField(record, "AlternatePosition1"),
			stringField(record, "AlternatePosition2"),
		}
	}
	return out, nil
}
