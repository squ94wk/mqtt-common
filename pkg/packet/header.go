package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

type header struct {
	pktType Type
	flags   byte
	length  uint32
}

func (h header) MsgType() Type {
	return h.pktType
}

func (h header) Flags() byte {
	return h.flags
}

func (h header) Length() uint32 {
	return h.length
}

func readHeader(reader io.Reader, header *header) error {
	var buf [1]byte
	if _, err := io.ReadFull(reader, buf[:1]); err != nil {
		return fmt.Errorf("failed to read packet header: %v", err)
	}

	var remainingLength uint32
	remainingLength, err := types.ReadVarInt(reader)
	if err != nil {
		return fmt.Errorf("failed to read packet header: %v", err)
	}

	header.flags = buf[0] >> 4
	header.pktType = Type(int8(buf[0] >> 4 << 4))
	header.length = uint32(remainingLength)
	return nil
}
