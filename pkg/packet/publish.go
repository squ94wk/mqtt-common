package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
	"github.com/squ94wk/mqtt-common/pkg/topic"
)

//Publish defines the publish control packet.
type Publish struct {
	Dup      bool
	Qos      byte
	Retain   bool
	Topic    topic.Topic
	PacketID uint16
	Props    Properties
	Payload  []byte
}

//WriteTo writes the publish control packet to writer according to the mqtt protocol.
func (p Publish) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	n1, err := writeFixedPublishHeader(p, writer)
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write publish packet: failed to write fixed header: %v", err)
	}

	n2, err := writeVariablePublishHeader(p, writer)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write publish packet: failed to write variable header: %v", err)
	}

	n3, err := writer.Write(p.Payload)
	n += int64(n3)
	if err != nil {
		return n, fmt.Errorf("failed to write publish packet: failed to write payload: %v", err)
	}

	return n, nil
}

func writeFixedPublishHeader(p Publish, writer io.Writer) (int64, error) {
	var n int64
	firstHeaderByte := byte(PUBLISH) << 4
	if p.Retain {
		firstHeaderByte |= 1
	}
	firstHeaderByte |= p.Qos << 1
	if p.Dup {
		firstHeaderByte |= 1 << 3
	}
	n1, err := writer.Write([]byte{firstHeaderByte})
	n += int64(n1)
	if err != nil {
		return n, fmt.Errorf("failed to write first header byte: %v", err)
	}

	// Remaining length
	var remainingLength = types.StringSize(p.Topic.String())
	remainingLength += types.UInt16Size
	remainingLength += p.Props.size()
	remainingLength += uint32(len(p.Payload))
	n2, err := types.WriteVarIntTo(writer, remainingLength)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write remaining packet length: %v", err)
	}

	return n, nil
}

func writeVariablePublishHeader(p Publish, writer io.Writer) (int64, error) {
	var n int64
	// 3.3.2.1 Topic name
	n1, err := types.WriteStringTo(writer, p.Topic.String())
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write topic name: %v", err)
	}
	// 3.3.2.2 Packet ID
	n2, err := types.WriteUInt16To(writer, p.PacketID)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write packet ID: %v", err)
	}

	// 3.3.2.3 Properties
	n3, err := p.Props.WriteTo(writer)
	n += n3
	if err != nil {
		return n, fmt.Errorf("failed to write properties: %v", err)
	}

	return n, nil
}

func readPublish(reader *io.LimitedReader, publish *Publish, headerFirstByte byte) error {
	// 3.3.1 Fixed header
	publish.Retain = headerFirstByte&1 > 0
	qos := headerFirstByte & (3 << 1) >> 1
	if qos == 3 {
		return fmt.Errorf("malformed packet: invalid QoS value %d", qos)
	}
	publish.Qos = qos
	publish.Dup = headerFirstByte&(1<<3) > 0

	// 3.3.2 Variable header
	// 3.3.2.1 Topic Name
	topicName, err := types.ReadString(reader)
	if err != nil {
		return fmt.Errorf("failed to read Publish packet: failed to read topic name: %v", err)
	}
	parsedTopic, err := topic.ParseTopic(topicName)
	if err != nil {
		return fmt.Errorf("malformed packet: invalid topic name: %s", topicName)
	}
	publish.Topic = parsedTopic

	// 3.3.2.2 Packet ID
	packetID, err := types.ReadUInt16(reader)
	if err != nil {
		return fmt.Errorf("failed to read Publish packet: failed to read packet ID: %v", err)
	}
	if packetID == 0 {
		return fmt.Errorf("malformed packet: invalid packet ID: %d", packetID)
	}
	publish.PacketID = packetID

	// 3.3.2.3 Properties
	props, err := readProperties(reader)
	if err != nil {
		return fmt.Errorf("failed to read publish packet: failed to read properties: %v", err)
	}
	publish.Props = props

	// 3.3.3 Payload
	payload := make([]byte, reader.N)
	_, err = io.ReadFull(reader, payload)
	if err != nil {
		return fmt.Errorf("failed to read publish packet: failed to read payload: %v", err)
	}

	publish.Payload = payload
	return nil
}
