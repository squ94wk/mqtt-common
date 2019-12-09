package types

import (
	"fmt"
	"io"
)

//WriteUInt16To writes a 16 bit integer to writer.
func WriteUInt16To(writer io.Writer, value uint16) (int64, error) {
	encoded := []byte{
		byte(value >> 8),
		byte(value),
	}
	n, err := writer.Write(encoded)
	if err != nil {
		return int64(n), fmt.Errorf("failed to write uint16: %v", err)
	}

	return 2, nil
}

//ReadUInt16 reads a 16 bit integer from reader.
func ReadUInt16(reader io.Reader) (uint16, error) {
	var buf [2]byte
	_, err := io.ReadFull(reader, buf[:2])
	if err != nil {
		return 0, fmt.Errorf("failed to read uint16: %v", err)
	}

	var value uint16
	value = uint16(buf[1])
	value |= uint16(buf[0]) << 8
	return value, nil
}
