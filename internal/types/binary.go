package types

import (
	"fmt"
	"io"
)

const maxBinaryLength = 1<<16 - 1

func WriteBinary(writer io.Writer, value []byte) error {
	if len(value) > maxBinaryLength {
		return fmt.Errorf("failed to write binary type: value too long ('%d' bytes > max = %d)", len(value), 2<<15)
	}

	size := uint16(len(value))
	err := WriteUInt16(writer, size)
	if err != nil {
		return fmt.Errorf("failed to write binary type: failed to write size '%d' encoded as two byte integer %v", size, err)
	}

	_, err = writer.Write(value)
	if err != nil {
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
	_, err = io.ReadFull(reader, buf[:])
	if err != nil {
		return nil, fmt.Errorf("failed to read binary type: failed to read payload: %v", err)
	}

	return buf, nil
}
