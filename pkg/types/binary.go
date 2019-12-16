package types

import (
	"fmt"
	"io"
)

func WriteBinary(writer io.Writer, value []byte) error {
	size := len(value)
	if size > int(^uint16(0)) {
		return fmt.Errorf("failed to write binary type: value too long ('%d' bytes > max = %d)", size, 2<<15)
	}

	if err := WriteUInt16(writer, uint16(size)); err != nil {
		return fmt.Errorf("failed to write binary type: failed to write size '%d' encoded as two byte integer %v", size, err)
	}

	if _, err := writer.Write(value); err != nil {
		return fmt.Errorf("failed to write binary type: failed to write payload: %v", err)
	}

	return nil
}

func ReadBinary(reader io.Reader) ([]byte, error) {
	size, err := ReadUInt16(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read binary type: failed to read size as encoded two byte integer: %v", err)
	}

	if size == 0 {
		return []byte{}, nil
	}

	buf := make([]byte, size)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, fmt.Errorf("failed to read binary type: failed to read payload: %v", err)
	}

	return buf, nil
}
