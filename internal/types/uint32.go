package types

import (
	"fmt"
	"io"
)

func WriteUInt32(writer io.Writer, value uint32) error {
	encoded := []byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	}

	_, err := writer.Write(encoded)
	if err != nil {
		return fmt.Errorf("failed to write uint32: %v", err)
	}
	return nil
}

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
