package dynasty

import (
	"os"
	"testing"
)

func TestDetectFormatFrTk(t *testing.T) {
	data := []byte("FrTk" + string(make([]byte, 64)))
	got := DetectFormat(data)
	if got.Compressed {
		t.Fatal("expected uncompressed")
	}
	if got.Format != FormatUncompressed {
		t.Fatalf("format = %q, want %q", got.Format, FormatUncompressed)
	}
}

func TestDetectFormatZlib(t *testing.T) {
	data := []byte{0x78, 0x9c, 0x01, 0x02, 0x03}
	got := DetectFormat(data)
	if !got.Compressed {
		t.Fatal("expected compressed")
	}
	if got.Format != FormatDynastyCommon {
		t.Fatalf("format = %q, want %q", got.Format, FormatDynastyCommon)
	}
}

func TestOpenEmptyFile(t *testing.T) {
	path := t.TempDir() + "/empty.sav"
	if err := writeFile(path, []byte{}); err != nil {
		t.Fatal(err)
	}
	_, err := Open(path, nil)
	if err == nil {
		t.Fatal("expected error for empty file")
	}
}

func TestInspectSyntheticSave(t *testing.T) {
	path := t.TempDir() + "/test.sav"
	payload := append([]byte("FrTk"), make([]byte, 60)...)
	if err := writeFile(path, payload); err != nil {
		t.Fatal(err)
	}

	file, err := Open(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	info, err := file.Inspect()
	if err != nil {
		t.Fatal(err)
	}
	if info.Size != int64(len(payload)) {
		t.Fatalf("size = %d, want %d", info.Size, len(payload))
	}
	if info.Compressed {
		t.Fatal("expected uncompressed synthetic save")
	}
}

func TestParseNotImplemented(t *testing.T) {
	path := t.TempDir() + "/test.sav"
	payload := append([]byte("FrTk"), make([]byte, 60)...)
	if err := writeFile(path, payload); err != nil {
		t.Fatal(err)
	}

	file, err := Open(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Parse(); !IsNotImplemented(err) {
		t.Fatalf("Parse() = %v, want ErrNotImplemented", err)
	}
}

func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o644)
}
