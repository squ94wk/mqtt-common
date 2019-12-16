package types

import (
	"fmt"
	"io"
)

func WriteUInt16(writer io.Writer, value uint16) error {
	encoded := []byte{
		byte(value >> 8),
		byte(value),
	}
	if _, err := writer.Write(encoded); err != nil {
		return fmt.Errorf("failed to write uint16: %v", err)
	}

	return nil
}

func ReadUInt16(reader io.Reader) (uint16, error) {
	var buf [2]byte
	if _, err := io.ReadFull(reader, buf[:2]); err != nil {
		return 0, fmt.Errorf("failed to read uint16: %v", err)
	}

	var value uint16
	value = uint16(buf[1])
	value |= uint16(buf[0]) << 8
	return value, nil
}
