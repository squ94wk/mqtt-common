package types

import (
	"fmt"
	"io"
)

const utf8StringMaxLength int = 1<<16 - 1

//WriteStringTo writes a string to writer.
func WriteStringTo(writer io.Writer, value string) (int64, error) {
	if len(value) > utf8StringMaxLength {
		return 0, fmt.Errorf("length of string exceeds maximum allowed length of %d bytes", utf8StringMaxLength)
	}

	var n int64
	length := uint16(len(value))
	n1, err := WriteUInt16To(writer, length)
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write length of UTF8 encoded string: %v", err)
	}

	n2, err := writer.Write([]byte(value))
	n += int64(n2)
	if err != nil {
		return n, fmt.Errorf("failed to write UTF8 encoded string: %v", err)
	}

	return n, nil
}

//ReadString reads a string from reader.
func ReadString(reader io.Reader) (string, error) {
	size, err := ReadUInt16(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read length of UTF8 encoded string: %v", err)
	}

	if size == 0 {
		return "", nil
	}

	buf := make([]byte, size)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return "", fmt.Errorf("failed to read UTF8 encoded string: %v", err)
	}

	return string(buf), nil
}
