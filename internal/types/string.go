package types

import (
	"fmt"
	"io"
)

const UTF8StringMaxLength uint16 = 1<<16 - 1

func WriteString(writer io.Writer, value string) error {
	length := uint16(len(value))
	err := WriteUInt16(writer, length)
	if err != nil {
		return fmt.Errorf("failed to write length of UTF8 encoded string: %v", err)
	}

	_, err = writer.Write([]byte(value))
	if err != nil {
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
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return "", fmt.Errorf("failed to read UTF8 encoded string: %v", err)
	}

	return string(buf), nil
}
