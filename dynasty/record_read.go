package dynasty

import "fmt"

// ReadRecords decodes all rows for a table using its attached schema.
func (t *Table) ReadRecords() error {
	if t == nil {
		return fmt.Errorf("cfb-dynasty: nil table")
	}
	if len(t.Records) > 0 {
		return nil
	}
	if t.Schema == nil {
		return fmt.Errorf("cfb-dynasty: table %q: schema required to read records", t.Name())
	}
	if len(t.Data) == 0 {
		return fmt.Errorf("cfb-dynasty: table %q: no data", t.Name())
	}

	offsets, err := buildOffsetTable(t.Data, t.Schema, t.Header)
	if err != nil {
		return err
	}

	table2 := t.Data
	if t.Header.HasSecondTable && t.Header.Table2StartIndex < len(t.Data) {
		table2 = t.Data[t.Header.Table2StartIndex:]
	}

	records := make([]Record, 0, t.Header.RecordCount)
	for i := 0; i < int(t.Header.RecordCount); i++ {
		start := t.Header.Table1StartIndex + i*t.Header.RecordSize
		end := start + t.Header.RecordSize
		if end > len(t.Data) {
			break
		}
		row := t.Data[start:end]
		fields := make(map[string]FieldValue, len(offsets))
		for _, entry := range offsets {
			fields[entry.Name] = decodeFieldValue(row, table2, entry, offsets)
		}
		records = append(records, Record{
			Index:  i,
			Fields: fields,
		})
	}
	t.Records = records
	return nil
}

// ReadAllRecords decodes rows for every table that has an attached schema.
func (f *File) ReadAllRecords() error {
	if !f.loaded {
		return ErrNotLoaded
	}
	for i := range f.tables {
		if f.tables[i].Schema == nil {
			continue
		}
		if err := f.tables[i].ReadRecords(); err != nil {
			return err
		}
	}
	return nil
}

// ReadTableRecords decodes rows for the first table with the given name.
func (f *File) ReadTableRecords(name string) error {
	if !f.loaded {
		return ErrNotLoaded
	}
	table, ok := f.GetTableByName(name)
	if !ok {
		return fmt.Errorf("%w: %q", ErrTableNotFound, name)
	}
	return table.ReadRecords()
}

// PrimaryTableByName returns the table instance with the most allocated rows.
func (f *File) PrimaryTableByName(name string) (*Table, bool) {
	if !f.loaded {
		return nil, false
	}
	var best *Table
	for i := range f.tables {
		table := &f.tables[i]
		if table.Name() != name {
			continue
		}
		if best == nil || table.AllocatedRecordCount() > best.AllocatedRecordCount() {
			best = table
		}
	}
	if best == nil {
		return nil, false
	}
	return best, true
}
