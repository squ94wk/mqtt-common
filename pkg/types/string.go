package types

import (
	"fmt"
	"io"
)

const UTF8StringMaxLength = 65535

func WriteString(writer io.Writer, value string) error {
	if err := WriteUInt16(writer, uint16(len(value))); err != nil {
		return fmt.Errorf("failed to write length of UTF8 encoded string: %v", err)
	}
	if _, err := writer.Write([]byte(value)); err != nil {
		return fmt.Errorf("failed to write UTF8 encoded string: %v", err)
	}

	return nil
}

func ReadString(reader io.Reader) (string, error) {
	size, err := ReadUInt16(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read length of UTF8 encoded string: %v", err)
	}

	if size == 0 {
		return "", nil
	}

	buf := make([]byte, size)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return "", fmt.Errorf("failed to read UTF8 encoded string: %v", err)
	}

	return string(buf), nil
}
