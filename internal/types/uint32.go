package types

import (
	"fmt"
	"io"
)

//WriteUInt32To writes a 32 bit integer to writer.
func WriteUInt32To(writer io.Writer, value uint32) (int64, error) {
	encoded := []byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	}

	n, err := writer.Write(encoded)
	if err != nil {
		return int64(n), fmt.Errorf("failed to write uint32: %v", err)
	}
	return 4, nil
}

//ReadUInt32 reads a 32 bit integer from reader.
func ReadUInt32(reader io.Reader) (uint32, error) {
	var buf [4]byte
	_, err := io.ReadFull(reader, buf[:])
	if err != nil {
		return 0, fmt.Errorf("failed to read uint32: %v", err)
	}

	var value uint32
	value = uint32(buf[3])
	value |= uint32(buf[2]) << 8
	value |= uint32(buf[1]) << 16
	value |= uint32(buf[0]) << 24
	return value, nil
}
