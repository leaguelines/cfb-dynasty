package binary

import "testing"

func TestHasPrefix(t *testing.T) {
	data := []byte("FrTkSPBF")
	if !HasPrefix(data, []byte("FrTk")) {
		t.Fatal("expected FrTk prefix")
	}
	if HasPrefix(data, []byte("SPBF")) {
		t.Fatal("did not expect SPBF at start")
	}
}

func TestFindAll(t *testing.T) {
	data := []byte("aaaSPBFbbbSPBFccc")
	got := FindAll(data, []byte("SPBF"))
	if len(got) != 2 || got[0] != 3 || got[1] != 10 {
		t.Fatalf("FindAll() = %v, want [3 10]", got)
	}
}

func TestHexPrefix(t *testing.T) {
	got := HexPrefix([]byte{0x46, 0x72, 0x54, 0x6b}, 4)
	want := "4672546b"
	if got != want {
		t.Fatalf("HexPrefix() = %q, want %q", got, want)
	}
}
