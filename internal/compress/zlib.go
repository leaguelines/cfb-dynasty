// Package compress unpacks dynasty save containers.
package compress

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

// Container describes how compressed payload bytes are laid out.
type Container struct {
	Compressed bool
	// Skip is the byte offset before zlib data when Compressed is true.
	Skip int
}

// Unpack returns decompressed bytes. When Compressed is false, data is returned as-is.
func Unpack(data []byte, container Container) ([]byte, error) {
	if !container.Compressed {
		return data, nil
	}
	return inflate(data, container.Skip)
}

func inflate(data []byte, offset int) ([]byte, error) {
	if offset >= len(data) {
		return nil, fmt.Errorf("cfb-dynasty: inflate offset %d out of range (len=%d)", offset, len(data))
	}
	reader, err := zlib.NewReader(bytes.NewReader(data[offset:]))
	if err != nil {
		return nil, fmt.Errorf("cfb-dynasty: zlib reader: %w", err)
	}
	defer reader.Close()

	out, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("cfb-dynasty: zlib inflate: %w", err)
	}
	return out, nil
}
