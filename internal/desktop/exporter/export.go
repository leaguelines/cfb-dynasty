package exporter

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"github.com/leaguelines/cfb-dynasty/dynasty"
)

// FullJSONStream writes pretty-printed JSON for the entire export directly to w,
// avoiding buffering the whole payload in memory.
func FullJSONStream(w io.Writer, e dynasty.Export) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(e)
}

// SectionJSONStream writes pretty-printed JSON for a single collection to w.
func SectionJSONStream(w io.Writer, e dynasty.Export, c Collection) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(c.JSON(e))
}

// SectionCSV writes a single collection's CSV to w.
func SectionCSV(w io.Writer, e dynasty.Export, c Collection) error {
	return WriteCSV(w, c.Header, c.Rows(e))
}

// WriteCSV writes a header and rows as CSV to w.
func WriteCSV(w io.Writer, header []string, rows [][]string) error {
	cw := csv.NewWriter(w)
	if err := cw.Write(header); err != nil {
		return err
	}
	if err := cw.WriteAll(rows); err != nil {
		return err
	}
	cw.Flush()
	return cw.Error()
}

// FullCSVZip writes a zip archive containing one CSV per non-empty collection.
func FullCSVZip(w io.Writer, e dynasty.Export) error {
	zw := zip.NewWriter(w)
	for _, c := range Registry {
		if c.Count(e) == 0 {
			continue
		}
		f, err := zw.Create(c.Name + ".csv")
		if err != nil {
			return err
		}
		if err := SectionCSV(f, e, c); err != nil {
			return err
		}
	}
	if err := zw.Close(); err != nil {
		return fmt.Errorf("finalize zip: %w", err)
	}
	return nil
}
