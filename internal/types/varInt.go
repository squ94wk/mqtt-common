package types

import (
	"fmt"
	"io"
)

const (
	b01111111 = 1<<7 - 1
	b10000000 = 1 << 7
)

//WriteVarIntTo writes an encoded variable length integer to writer.
func WriteVarIntTo(writer io.Writer, value uint32) (int64, error) {
	var buf [4]byte
	encoded := encodeVarInt(value, buf)

	n, err := writer.Write(encoded)
	if err != nil {
		return int64(n), fmt.Errorf("couldn't write varInt '%d': %v", value, err)
	}

	return int64(n), nil
}

//ReadVarInt reads an encoded variable length integer from reader.
func ReadVarInt(reader io.Reader) (uint32, error) {
	var offset, value uint32 = 0, 0
	var buf [1]byte
	for pos := 0; pos < 4; pos++ {
		length, err := reader.Read(buf[:])
		if err != nil || length == 0 {
			return 0, fmt.Errorf("failed to read byte (current value: %d): %v", value, err)
		}

		value += uint32(buf[0]&b01111111) << offset

		offset += 7
		if (buf[0] & b10000000) == 0 {
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
