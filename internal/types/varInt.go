package types

import (
	"fmt"
	"io"
)

const (
	b0111_1111 = 1<<7 - 1
	b1000_0000 = 1 << 7
)

func WriteVarInt(writer io.Writer, value uint32) error {
	var buf [4]byte
	encoded := encodeVarInt(value, buf)

	_, err := writer.Write(encoded)
	if err != nil {
		return fmt.Errorf("couldn't write varInt '%d': %v", value, err)
	}

	return nil
}

func ReadVarInt(reader io.Reader) (uint32, error) {
	var offset, value uint32 = 0, 0
	var buf [1]byte
	for pos := 0; pos < 4; pos++ {
		length, err := reader.Read(buf[:])
		if err != nil || length == 0 {
			return 0, fmt.Errorf("failed to read byte (current value: %d): %v", value, err)
		}

		value += uint32(buf[0]&b0111_1111) << offset

		offset += 7
		if (buf[0] & b1000_0000) == 0 {
			return value, nil
		}
	}

	return 0, fmt.Errorf("malformed varint: value would exceed maximum")
}

func encodeVarInt(varInt uint32, buf [4]byte) []byte {
	if varInt == 0 {
		return []byte{0}
	}

	pos, value := 0, varInt
	for ; value > 0 && pos < 4; pos++ {
		encodedByte := value % 128
		value = value / 128
		if value > 0 {
			// set MSB = 1
			encodedByte = encodedByte | 128
		}
		buf[pos] = byte(encodedByte)
	}

	if value > 0 {
		panic(fmt.Sprintf("can't encode VarInt: value '%d' exceeds maximum", varInt))
	}

	return buf[:pos]
}
