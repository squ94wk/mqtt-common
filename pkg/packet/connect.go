package packet

import (
	"bytes"
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
	"github.com/squ94wk/mqtt-common/pkg/topic"
)

type Connect struct {
	keepAlive  uint16
	cleanStart bool
	properties
	payload ConnectPayload
}

type ConnectPayload struct {
	clientId string
	properties
	willRetain  bool
	willQoS     QoS
	willTopic   topic.Topic
	willPayload []byte
	username    string
	password    []byte
}

func NewConnect(
	cleanStart bool,
	keepAlive uint16,
	props map[PropId][]Property,
	clientId string,
	willRetain bool,
	willQoS QoS,
	willProps map[PropId][]Property,
	willTopic topic.Topic,
	willPayload []byte,
	username string,
	password []byte) Connect {
	return Connect{
		keepAlive:  keepAlive,
		cleanStart: cleanStart,
		properties: properties{props},
		payload: ConnectPayload{
			clientId:    clientId,
			properties:  properties{willProps},
			willRetain:  willRetain,
			willQoS:     willQoS,
			willTopic:   willTopic,
			willPayload: willPayload,
			username:    username,
			password:    password,
		},
	}
}

func (c Connect) KeepAlive() uint16 {
	return c.keepAlive
}
func (c *Connect) SetKeepAlive(value uint16) {
	c.keepAlive = value
}

func (c Connect) CleanStart() bool {
	return c.cleanStart
}
func (c *Connect) SetCleanStart(cleanStart bool) {
	c.cleanStart = cleanStart
}

func (c Connect) Props() map[PropId][]Property {
	return c.props
}

func (c Connect) Payload() ConnectPayload {
	return c.payload
}

func (p ConnectPayload) ClientId() string {
	return p.clientId
}
func (p *ConnectPayload) SetClientId(clientId string) {
	p.clientId = clientId
}

func (p ConnectPayload) WillProps() map[PropId][]Property {
	return p.props
}

func (p ConnectPayload) WillTopic() topic.Topic {
	return p.willTopic
}
func (p *ConnectPayload) SetWillTopic(topic topic.Topic) {
	p.willTopic = topic
}

func (p ConnectPayload) WillQoS() QoS {
	return p.willQoS
}
func (p *ConnectPayload) SetWillQoS(qos QoS) {
	p.willQoS = qos
}

func (p ConnectPayload) WillPayload() []byte {
	return p.willPayload
}
func (p *ConnectPayload) SetWillPayload(payload []byte) {
	p.willPayload = payload
}

func (p ConnectPayload) Username() string {
	return p.username
}
func (p *ConnectPayload) SetUsername(username string) {
	p.username = username
}

func (p ConnectPayload) Password() []byte {
	return p.password
}
func (p *ConnectPayload) SetPassword(password []byte) {
	p.password = password
}

func (p ConnectPayload) WillRetain() bool {
	return p.willRetain
}
func (p *ConnectPayload) SetWillRetain(willRetain bool) {
	p.willRetain = willRetain
}

func ReadConnect(origReader io.Reader, connect *Connect, header header) error {
	reader := io.LimitReader(origReader, int64(header.Length()))

	// 3.1.2 Variable header
	var buf [8]byte
	if _, err := io.ReadFull(reader, buf[:8]); err != nil {
		return fmt.Errorf("failed to read connect packet: failed to read variable header: %v", err)
	}

	// 3.1.2.1 Protocol Name
	if bytes.Compare(buf[:6], []byte{0, 4, 'M', 'Q', 'T', 'T'}) != 0 {
		return fmt.Errorf("unsupported protocol name '%v'", buf[:6])
	}

	// 3.1.2.2 Protocol Version
	if buf[6] != 5 {
		return fmt.Errorf("unsupported protocol version '%d'", buf[6])
	}
	// protocol version is not saved, since only '5' is supported

	// 3.1.2.3 Connect Flags
	flags := buf[7]
	reserved := byte(flags)&byte(1) >= 1
	cleanStart := byte(flags)&byte(1<<1) >= 1
	hasWill := byte(flags)&byte(1<<2) >= 1
	willQoS := QoS(flags & 24 >> 3)
	willRetain := (flags)&byte(1<<5) >= 1
	hasPassword := byte(flags)&byte(1<<6) >= 1
	hasUsername := byte(flags)&byte(1<<7) >= 1

	if reserved {
		return fmt.Errorf("malformed packet: reserved flag is set")
	}
	if !hasWill {
		if willRetain {
			return fmt.Errorf("malformed packet: will retain flag is set, but will flag is not set")
		}
		if willQoS > 0 {
			return fmt.Errorf("malformed packet: will QoS is > 0, but will flag is not set")
		}
	}

	connect.SetCleanStart(cleanStart)
	connect.payload.SetWillRetain(willRetain)

	// 3.1.2.6 Will QoS
	if willQoS > 2 {
		return fmt.Errorf("malformed packet: Will QoS must be 0, 1 or 2, but is %d", willQoS)
	}

	// 3.1.2.10 Keep Alive
	keepAlive, err := types.ReadUInt16(reader)
	if err != nil {
		return fmt.Errorf("failed to read connect packet: failed to read keepAlive: %v", err)
	}
	connect.SetKeepAlive(keepAlive)

	// 3.1.2.11 Properties
	connect.ResetProps()
	if err := readProperties(reader, connect.props); err != nil {
		return fmt.Errorf("failed to read connect packet: failed to read will properties: %v", err)
	}

	// 3.1.3 Payload
	payload := ConnectPayload{}
	// 3.1.3.1 ClientID
	clientId, err := types.ReadString(reader)
	if err != nil {
		return fmt.Errorf("failed to read clientID: %v", err)
	}
	clientIdLength := len(clientId)
	if clientIdLength > 23 {
		return fmt.Errorf("malformed packet: ClientID too long (%d), must be between 1 and 23 char long", clientIdLength)
	}
	//TODO: check characters are only in a-zA-Z0-9
	payload.SetClientId(clientId)

	if hasWill {
		payload.SetWillQoS(willQoS)

		payload.ResetProps()
		// 3.1.3.2 Will properties
		err := readProperties(reader, payload.props)
		if err != nil {
			return fmt.Errorf("failed to read connect packet: failed to read will properties: %v", err)
		}

		// 3.1.3.3 Will topic
		willTopic, err := types.ReadString(reader)
		if err != nil {
			return fmt.Errorf("failed to read will topic: %v", err)
		}
		payload.SetWillTopic(topic.Topic(willTopic))

		// 3.1.3.4 Will payload
		willPayload, err := types.ReadBinary(reader)
		if err != nil {
			return fmt.Errorf("failed to read will payload: %v", err)
		}
		payload.SetWillPayload(willPayload)
	}

	// 3.1.3.5 User name
	if hasUsername {
		username, err := types.ReadString(reader)
		if err != nil {
			return fmt.Errorf("failed to read username: %v", err)
		}
		payload.SetUsername(username)
	}

	// 3.1.3.6 Password
	if hasPassword {
		password, err := types.ReadBinary(reader)
		if err != nil {
			return fmt.Errorf("failed to read password: %v", err)
		}
		payload.SetPassword(password)
	}

	connect.payload = payload

	return nil
}

func (c Connect) Write(writer io.Writer) error {
	// 3.1.1 Fixed header
	firstHeaderByte := byte(CONNECT) << 4
	if _, err := writer.Write([]byte{firstHeaderByte}); err != nil {
		return fmt.Errorf("failed to write connect packet: failed to write fixed header: %v", err)
	}

	// Remaining length
	var remainingLength uint32 = 7 + 1 + 2 // protocol name & version, flags, keep alive
	remainingLength += propertiesSize(c.props)
	remainingLength += types.StringSize(c.payload.clientId)
	if c.payload.willTopic != "" {
		remainingLength += propertiesSize(c.payload.props)
		remainingLength += types.StringSize(string(c.payload.willTopic))
		remainingLength += types.BinarySize(c.payload.willPayload)
	}
	if c.payload.username != "" {
		remainingLength += types.StringSize(c.payload.username)
	}
	if c.payload.password != nil {
		remainingLength += types.BinarySize(c.payload.password)
	}
	if err := types.WriteVarInt(writer, remainingLength); err != nil {
		return fmt.Errorf("failed to write connect packet: failed to write remaining packet length: %v", err)
	}

	// 3.1.2 Variable header
	// 3.1.2.1 Protocol Name
	// 3.1.2.2 Protocol Version
	if _, err := writer.Write([]byte{0, 4, 'M', 'Q', 'T', 'T', 5}); err != nil {
		return fmt.Errorf("failed to write connect packet: failed to write protocol info: %v", err)
	}

	// 3.1.2.3 Connect Flags
	var flags byte = 0
	// reserved
	flags |= 0
	if c.cleanStart {
		flags |= 1 << 1
	}
	if c.payload.willTopic != "" {
		flags |= 1 << 2
	}
	switch c.payload.willQoS {
	case Qos0:
	case Qos1:
		flags |= 1 << 3
	case Qos2:
		flags |= 1 << 4
	default:
		panic(fmt.Errorf("invalid option for QoS: %v", c.payload.willQoS))
	}
	if c.payload.willRetain {
		flags |= 1 << 5
	}
	if c.payload.password != nil {
		flags |= 1 << 6
	}
	if c.payload.username != "" {
		flags |= 1 << 7
	}

	if _, err := writer.Write([]byte{flags}); err != nil {
		return fmt.Errorf("failed to write connect packet: failed to write flags: %v", err)
	}

	// 3.1.2.10 Keep Alive
	if err := types.WriteUInt16(writer, c.keepAlive); err != nil {
		return fmt.Errorf("failed to write connect packet: failed to write keep alive: %v", err)
	}

	// 3.1.2.11 Properties
	if err := WriteProperties(writer, c.props); err != nil {
		return fmt.Errorf("failed to write connect packet: %v", err)
	}

	// 3.1.3 Payload
	// 3.1.3.1 ClientID
	if err := types.WriteString(writer, c.payload.clientId); err != nil {
		return fmt.Errorf("failed to write connect packet: failed to write client id: %v", err)
	}

	// 3.1.3.2 Will properties
	if c.payload.willTopic != "" {
		// 3.1.3.2.1 Property length
		if err := WriteProperties(writer, c.payload.props); err != nil {
			return fmt.Errorf("failed to write connect packet: failed to write will properties: %v", err)
		}

		// 3.1.3.3 Will topic
		if err := types.WriteString(writer, string(c.payload.willTopic)); err != nil {
			return fmt.Errorf("failed to write connect packet: failed to write will topic: %v", err)
		}

		// 3.1.3.4 Will payload
		if err := types.WriteBinary(writer, c.payload.willPayload); err != nil {
			return fmt.Errorf("failed to write connect packet: failed to write will payload: %v", err)
		}
	}

	// 3.1.3.5 User name
	if c.payload.username != "" {
		if err := types.WriteString(writer, c.payload.username); err != nil {
			return fmt.Errorf("failed to write connect packet: failed to write username: %v", err)
		}
	}

	// 3.1.3.6 Password
	if c.payload.password != nil {
		if err := types.WriteBinary(writer, c.payload.password); err != nil {
			return fmt.Errorf("failed to write Connect packet: failed to write password: %v", err)
		}
	}

	return nil
}
