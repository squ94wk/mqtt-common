package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

//Suback defines the suback control packet.
type Suback struct {
	packetID uint16
	props    Properties
	reasons  []SubackReason
}

//PacketID returns the value of the suback control packet.
func (s Suback) PacketID() uint16 {
	return s.packetID
}

//SetPacketID sets the value of the suback control packet.
func (s *Suback) SetPacketID(packetID uint16) {
	s.packetID = packetID
}

//Props returns the properties of the suback control packet.
func (s Suback) Props() Properties {
	return s.props
}

//SetProps replaces the properties of the suback control packet.
func (s *Suback) SetProps(props map[uint32][]Property) {
	s.props = props
}

//Reasons returns to each subscription the corresponding highest supported QoS or an error reason code.
func (s Suback) Reasons() []SubackReason {
	return s.reasons
}

//SetReasons sets the reason codes for all subscriptions.
func (s *Suback) SetReasons(reasons []SubackReason) {
	s.reasons = reasons
}

//WriteTo writes the suback control packet to writer according to the mqtt protocol.
func (s Suback) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	// 3.8.1 Fixed header
	firstByte := byte(SUBACK) << 4
	n1, err := writer.Write([]byte{firstByte})
	n += int64(n1)
	if err != nil {
		return n, fmt.Errorf("failed to write suback packet: failed to write fixed header: %v", err)
	}

	//3.8.2 Variable header
	var remainingLength = types.UInt16Size // packetID
	remainingLength += s.props.size()
	remainingLength += uint32(len(s.reasons))
	n2, err := types.WriteVarIntTo(writer, remainingLength)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write suback packet: failed to write packet length: %v", err)
	}

	n4, err := types.WriteUInt16To(writer, s.packetID)
	n += n4
	if err != nil {
		return n, fmt.Errorf("failed to write suback packet: failed to write packetID: %v", err)
	}

	n5, err := s.props.WriteTo(writer)
	n += n5
	if err != nil {
		return n, fmt.Errorf("failed to write suback packet: failed to write properties: %v", err)
	}

	reasonsBuf := make([]byte, len(s.reasons))
	for i, reason := range s.reasons {
		reasonsBuf[i] = byte(reason)
	}
	n6, err := writer.Write(reasonsBuf)
	n += int64(n6)
	if err != nil {
		return n, fmt.Errorf("failed to write suback packet: failed to write suback reason codes: %v", err)
	}

	return n, nil
}

func readSuback(reader io.Reader, suback *Suback) error {
	// 3.8.2 Variable header
	// 3.8.2.1 Suback packet ID
	packetID, err := types.ReadUInt16(reader)
	if err != nil {
		return fmt.Errorf("failed to read suback packet: failed to read packet ID: %v", err)
	}
	suback.SetPacketID(packetID)

	// 3.8.2.2 Suback properties
	props, err := readProperties(reader)
	if err != nil {
		return fmt.Errorf("failed to read suback packet: failed to read properties: %v", err)
	}
	suback.props = props

	// 3.8.3 Payload
	reasonBuf := make([]byte, reader.(*io.LimitedReader).N)
	_, err = io.ReadFull(reader, reasonBuf)
	if err != nil {
		return fmt.Errorf("failed to read suback packet: failed to read reason codes: %v", err)
	}

	reasons := make([]SubackReason, len(reasonBuf))
	for i, reason := range reasonBuf {
		reasons[i] = SubackReason(reason)
	}
	suback.SetReasons(reasons)
	return nil
}
