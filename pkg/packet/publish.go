package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
	"github.com/squ94wk/mqtt-common/pkg/topic"
)

//Publish defines the publish control packet.
type Publish struct {
	dup      bool
	qos      byte
	retain   bool
	topic    topic.Topic
	packetID uint16
	props    Properties
	payload  []byte
}

//Dup returns the value of the DUP flag of the publish control packet.
func (p Publish) Dup() bool {
	return p.dup
}

//SetDup sets the value of the DUP flag of the publish control packet.
func (p *Publish) SetDup(dup bool) {
	p.dup = dup
}

//QoS returns the value of the QoS flag of the publish control packet.
func (p Publish) QoS() byte {
	return p.qos
}

//SetQoS sets the value of the QoS flag of the publish control packet.
func (p *Publish) SetQoS(qos byte) {
	p.qos = qos
}

//Retain returns the value of the Retain flag of the publish control packet.
func (p Publish) Retain() bool {
	return p.retain
}

//SetRetain sets the value of the Retain flag of the publish control packet.
func (p *Publish) SetRetain(retain bool) {
	p.retain = retain
}

//Topic returns the topic of the publish control packet.
func (p Publish) Topic() topic.Topic {
	return p.topic
}

//SetTopic sets the topic of the publish control packet.
func (p *Publish) SetTopic(topic topic.Topic) {
	p.topic = topic
}

//PacketID returns the value of the publish control packet.
func (p Publish) PacketID() uint16 {
	return p.packetID
}

//SetPacketID sets the value of the publish control packet.
func (p *Publish) SetPacketID(packetID uint16) {
	p.packetID = packetID
}

//Payload returns the payload of the publish control packet.
func (p Publish) Payload() []byte {
	return p.payload
}

//SetPayload sets the payload of the publish control packet.
func (p *Publish) SetPayload(payload []byte) {
	p.payload = payload
}

//Props returns the properties of the publish control packet.
func (p Publish) Props() Properties {
	return p.props
}

//SetProps replaces the properties of the publish control packet.
func (p *Publish) SetProps(props Properties) {
	p.props = props
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

	n3, err := writer.Write(p.payload)
	n += int64(n3)
	if err != nil {
		return n, fmt.Errorf("failed to write publish packet: failed to write payload: %v", err)
	}

	return n, nil
}

func writeFixedPublishHeader(p Publish, writer io.Writer) (int64, error) {
	var n int64
	firstHeaderByte := byte(PUBLISH) << 4
	if p.retain {
		firstHeaderByte |= 1
	}
	firstHeaderByte |= p.qos << 1
	if p.dup {
		firstHeaderByte |= 1 << 3
	}
	n1, err := writer.Write([]byte{firstHeaderByte})
	n += int64(n1)
	if err != nil {
		return n, fmt.Errorf("failed to write first header byte: %v", err)
	}

	// Remaining length
	var remainingLength = types.StringSize(p.topic.String())
	remainingLength += types.UInt16Size
	remainingLength += p.props.size()
	remainingLength += uint32(len(p.payload))
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
	n1, err := types.WriteStringTo(writer, p.topic.String())
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write topic name: %v", err)
	}
	// 3.3.2.2 Packet ID
	n2, err := types.WriteUInt16To(writer, p.packetID)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write packet ID: %v", err)
	}

	// 3.3.2.3 Properties
	n3, err := p.props.WriteTo(writer)
	n += n3
	if err != nil {
		return n, fmt.Errorf("failed to write properties: %v", err)
	}

	return n, nil
}

func readPublish(reader *io.LimitedReader, publish *Publish, headerFirstByte byte) error {
	// 3.3.1 Fixed header
	publish.retain = headerFirstByte&1 > 0
	qos := headerFirstByte & (3 << 1) >> 1
	if qos == 3 {
		return fmt.Errorf("malformed packet: invalid QoS value %d", qos)
	}
	publish.qos = qos
	publish.dup = headerFirstByte&(1<<3) > 0

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
	publish.topic = parsedTopic

	// 3.3.2.2 Packet ID
	packetID, err := types.ReadUInt16(reader)
	if err != nil {
		return fmt.Errorf("failed to read Publish packet: failed to read packet ID: %v", err)
	}
	if packetID == 0 {
		return fmt.Errorf("malformed packet: invalid packet ID: %d", packetID)
	}
	publish.packetID = packetID

	// 3.3.2.3 Properties
	props, err := readProperties(reader)
	if err != nil {
		return fmt.Errorf("failed to read publish packet: failed to read properties: %v", err)
	}
	publish.props = props

	// 3.3.3 Payload
	payload := make([]byte, reader.N)
	_, err = io.ReadFull(reader, payload)
	if err != nil {
		return fmt.Errorf("failed to read publish packet: failed to read payload: %v", err)
	}

	publish.payload = payload
	return nil
}
