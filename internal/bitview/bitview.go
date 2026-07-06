// Package bitview reads big-endian bit fields from FranTk table records.
package bitview

import (
	"encoding/binary"
	"math"
)

// GetBits reads length bits starting at bitOffset using MSB-first ordering.
func GetBits(data []byte, bitOffset, length int) uint32 {
	if length <= 0 || length > 32 {
		return 0
	}
	var result uint32
	for i := 0; i < length; i++ {
		byteIndex := (bitOffset + i) / 8
		if byteIndex >= len(data) {
			break
		}
		bitIndex := 7 - ((bitOffset + i) % 8)
		bit := (data[byteIndex] >> bitIndex) & 1
		result = (result << 1) | uint32(bit)
	}
	return result
}

// GetFloat32 reads a 32-bit float stored at bitOffset.
func GetFloat32(data []byte, bitOffset int) float32 {
	bits := GetBits(data, bitOffset, 32)
	return math.Float32frombits(bits)
}

// GetUint32AtBit reads a big-endian uint32 aligned at bitOffset (must be byte-aligned).
func GetUint32AtBit(data []byte, bitOffset int) uint32 {
	byteOffset := bitOffset / 8
	if byteOffset+4 > len(data) {
		return 0
	}
	return binary.BigEndian.Uint32(data[byteOffset:])
}
