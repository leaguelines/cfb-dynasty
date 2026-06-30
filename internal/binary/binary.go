// Package binary provides low-level helpers for scanning dynasty save bytes.
package binary

import (
	"encoding/hex"
	"fmt"
)

// HasPrefix reports whether data begins with prefix.
func HasPrefix(data, prefix []byte) bool {
	if len(data) < len(prefix) {
		return false
	}
	for i := range prefix {
		if data[i] != prefix[i] {
			return false
		}
	}
	return true
}

// FindAll returns every index of needle in haystack.
func FindAll(haystack, needle []byte) []int {
	if len(needle) == 0 || len(haystack) < len(needle) {
		return nil
	}

	var indices []int
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if HasPrefix(haystack[i:], needle) {
			indices = append(indices, i)
			if len(needle) == 1 {
				continue
			}
			i += len(needle) - 1
		}
	}
	return indices
}

// HexPrefix returns a spaced hex dump of the first n bytes.
func HexPrefix(data []byte, n int) string {
	if n <= 0 || len(data) == 0 {
		return ""
	}
	if n > len(data) {
		n = len(data)
	}
	return fmt.Sprintf("%s", hex.EncodeToString(data[:n]))
}
