package types

import (
	"fmt"
	"io"
)

const maxBinaryLength = 1<<16 - 1

//WriteBinaryTo writes a binary property to writer.
func WriteBinaryTo(writer io.Writer, value []byte) (int64, error) {
	if len(value) > maxBinaryLength {
		return 0, fmt.Errorf("failed to write binary type: value too long ('%d' bytes > max = %d)", len(value), 2<<15)
	}

	var n int64
	size := uint16(len(value))
	n1, err := WriteUInt16To(writer, size)
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write binary type: failed to write size '%d' encoded as two byte integer %v", size, err)
	}

	n2, err := writer.Write(value)
	n += int64(n2)
	if err != nil {
		return n, fmt.Errorf("failed to write binary type: failed to write payload: %v", err)
	}

	return n, nil
}

//ReadBinary reads a byte array from reader.
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
