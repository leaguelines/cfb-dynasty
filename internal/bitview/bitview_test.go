package bitview

import "testing"

func TestGetBits(t *testing.T) {
	data := []byte{0xA5}
	if got := GetBits(data, 0, 4); got != 0xA {
		t.Fatalf("high nibble = %#x, want 0xA", got)
	}
	if got := GetBits(data, 4, 4); got != 0x5 {
		t.Fatalf("low nibble = %#x, want 0x5", got)
	}
}

func TestGetFloat32(t *testing.T) {
	data := []byte{0x40, 0x48, 0x00, 0x00}
	got := GetFloat32(data, 0)
	if got < 3.124 || got > 3.126 {
		t.Fatalf("float = %v, want ~3.125", got)
	}
}
