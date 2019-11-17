package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

type header struct {
	pktType pktType
	flags   byte
	length  uint32
}

func readHeader(reader io.Reader, header *header) error {
	var buf [1]byte
	if _, err := io.ReadFull(reader, buf[:]); err != nil {
		return fmt.Errorf("failed to read packet header: %v", err)
	}

	var remainingLength uint32
	remainingLength, err := types.ReadVarInt(reader)
	if err != nil {
		return fmt.Errorf("failed to read packet header: %v", err)
	}

	header.flags = buf[0] << 4 >> 4
	header.pktType = pktType(buf[0] >> 4)
	header.length = remainingLength
	return nil
}
