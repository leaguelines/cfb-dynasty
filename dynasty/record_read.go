package dynasty

import (
	"fmt"
	"strings"
)

// ReadRecords decodes all rows for a table using its attached schema.
func (t *Table) ReadRecords() error {
	if t == nil {
		return fmt.Errorf("cfb-dynasty: nil table")
	}
	if len(t.Records) > 0 {
		return nil
	}
	if len(t.Data) == 0 {
		return fmt.Errorf("cfb-dynasty: table %q: no data", t.Name())
	}

	// Array-store tables (e.g. "Team[]") hold a fixed-width list of 32-bit
	// elements and carry no named schema in the bundle. Decode them with a
	// generated schema of RecordSize/4 element fields, matching the reference
	// parser's generic-array handling.
	if t.Schema == nil {
		if isArrayTable(t) {
			return t.readArrayRecords()
		}
		return fmt.Errorf("cfb-dynasty: table %q: schema required to read records", t.Name())
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

// isArrayTable reports whether a table is an array store (a fixed-width list of
// 32-bit elements) that should be decoded with a generated element schema.
func isArrayTable(t *Table) bool {
	if t.Header.IsArray {
		return true
	}
	return strings.HasSuffix(t.Name(), "[]")
}

// arrayElementIsReference reports whether an array table's elements are record
// references (element type starts uppercase, e.g. "Team[]") rather than scalars.
func arrayElementIsReference(name string) bool {
	base := strings.TrimSuffix(name, "[]")
	if base == "" {
		return false
	}
	c := base[0]
	return c >= 'A' && c <= 'Z'
}

// readArrayRecords decodes an array-store table into records whose fields are
// named by element index ("0", "1", ...). Each element is a 32-bit slot.
func (t *Table) readArrayRecords() error {
	if t.Header.RecordSize < 4 {
		t.Records = nil
		return nil
	}
	fieldsPerRow := t.Header.RecordSize / 4
	isRef := arrayElementIsReference(t.Name())

	records := make([]Record, 0, t.Header.RecordCount)
	for i := 0; i < int(t.Header.RecordCount); i++ {
		start := t.Header.Table1StartIndex + i*t.Header.RecordSize
		end := start + t.Header.RecordSize
		if end > len(t.Data) {
			break
		}
		row := t.Data[start:end]
		fields := make(map[string]FieldValue, fieldsPerRow)
		for f := 0; f < fieldsPerRow; f++ {
			bit := f * 32
			entry := OffsetEntry{Offset: bit, Length: 32, IsReference: isRef, Type: "int"}
			fields[fmt.Sprintf("%d", f)] = decodeFieldValue(row, nil, entry, nil)
		}
		records = append(records, Record{Index: i, Fields: fields})
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
