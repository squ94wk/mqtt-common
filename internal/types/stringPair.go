package types

import (
	"fmt"
	"io"
)

type StringPair struct {
	key   string
	value string
}

func NewStringPair(key string, value string) StringPair {
	return StringPair{key, value}
}

func (p StringPair) Key() string {
	return p.key
}

func (p StringPair) Value() string {
	return p.value
}

func WriteStringPair(writer io.Writer, pair StringPair) error {
	err := WriteString(writer, pair.Key())
	if err != nil {
		return fmt.Errorf("failed to write key of string pair: %v", err)
	}

	err = WriteString(writer, pair.Value())
	if err != nil {
		return fmt.Errorf("failed to write value of string pair: %v", err)
	}

	return nil
}

func ReadStringPair(reader io.Reader) (StringPair, error) {
	key, err := ReadString(reader)
	if err != nil {
		return StringPair{}, fmt.Errorf("failed to read key of string pair: %v", err)
	}

	value, err := ReadString(reader)
	if err != nil {
		return StringPair{}, fmt.Errorf("failed to read value of string pair: %v", err)
	}

	return StringPair{key, value}, nil
}
