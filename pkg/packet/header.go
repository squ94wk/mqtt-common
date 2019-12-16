package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/pkg/types"
)

type Header struct {
	pktType Type
	flags   byte
	length  uint32
}

func (h Header) MsgType() Type {
	return h.pktType
}

func (h Header) Flags() byte {
	return h.flags
}

func (h Header) Length() uint32 {
	return h.length
}

func ReadHeader(reader io.Reader, header *Header) error {
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

//func (h Header) Write(writer io.Writer) error {
//typeSpecificByte := (byte(h.pktType) << 4) & h.flags
//if err := writer.Write([]byte{typeSpecificByte}); err != nil {
//return fmt.Errorf("failed to write fixed header. couldn't write type & flags. %v", err)
//}

//// remaining length
//if err = types.WriteVarInt(writer, h.Length()); err != nil {
//return fmt.Errorf("failed to write fixed header. couldn't write length. %v", err)
//}

//return nil
//}
