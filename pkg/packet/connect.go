package packet

import (
	"bytes"
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

//Connect defines the connect control packet.
type Connect struct {
	keepAlive  uint16
	cleanStart bool
	props      Properties
	payload    ConnectPayload
}

//ConnectPayload defines the payload of a connect control packet.
type ConnectPayload struct {
	clientID    string
	willProps   Properties
	willRetain  bool
	willQoS     byte
	willTopic   string
	willPayload []byte
	username    string
	password    []byte
}

//KeepAlive returns the keep alive duration of the connect control packet in seconds.
func (c Connect) KeepAlive() uint16 {
	return c.keepAlive
}

//SetKeepAlive sets the keep alive duration of the connect control packet in seconds.
func (c *Connect) SetKeepAlive(value uint16) {
	c.keepAlive = value
}

//CleanStart returns the value of the clean start flag of the connect control packet.
func (c Connect) CleanStart() bool {
	return c.cleanStart
}

//SetCleanStart sets the value of the clean start flag of the connect control packet.
func (c *Connect) SetCleanStart(cleanStart bool) {
	c.cleanStart = cleanStart
}

//Props returns the properties of the connect control packet.
func (c Connect) Props() Properties {
	return c.props
}

//SetProps replaces the properties of the connect control packet.
func (c *Connect) SetProps(props Properties) {
	c.props = props
}

//Payload returns the payload of the connect control packet.
func (c Connect) Payload() ConnectPayload {
	return c.payload
}

//ClientID returns the client ID of the connect control packet payload.
func (p ConnectPayload) ClientID() string {
	return p.clientID
}

//SetClientID sets the client ID of the connect control packet payload.
func (p *ConnectPayload) SetClientID(clientID string) {
	p.clientID = clientID
}

//WillProps returns the properties of the will message of the connect control packet.
func (p ConnectPayload) WillProps() Properties {
	return p.willProps
}

//SetWillProps replaces the properties of the will message of the connect control packet.
func (p *ConnectPayload) SetWillProps(props Properties) {
	p.willProps = props
}

//WillTopic returns the topic of the will message of the connect control packet.
func (p ConnectPayload) WillTopic() string {
	return p.willTopic
}

//SetWillTopic sets the topic of the will message of the connect control packet.
func (p *ConnectPayload) SetWillTopic(topic string) {
	p.willTopic = topic
}

//WillQoS returns the quality of service level of the will message of the connect control packet.
func (p ConnectPayload) WillQoS() byte {
	return p.willQoS
}

//SetWillQoS sets the quality of service level of the will message of the connect control packet.
func (p *ConnectPayload) SetWillQoS(qos byte) {
	p.willQoS = qos
}

//WillPayload returns the payload of the will message of the connect control packet.
func (p ConnectPayload) WillPayload() []byte {
	return p.willPayload
}

//SetWillPayload sets the payload of the will message of the connect control packet.
func (p *ConnectPayload) SetWillPayload(payload []byte) {
	p.willPayload = payload
}

//Username returns the user name of the connect control packet.
func (p ConnectPayload) Username() string {
	return p.username
}

//SetUsername sets the user name of the connect control packet.
func (p *ConnectPayload) SetUsername(username string) {
	p.username = username
}

//Password returns the password of the connect control packet.
func (p ConnectPayload) Password() []byte {
	return p.password
}

//SetPassword sets the password of the connect control packet.
func (p *ConnectPayload) SetPassword(password []byte) {
	p.password = password
}

//WillRetain returns the value of the will retain flag of the connect control packet.
func (p ConnectPayload) WillRetain() bool {
	return p.willRetain
}

//SetWillRetain sets the value of the will retain flag of the connect control packet.
func (p *ConnectPayload) SetWillRetain(willRetain bool) {
	p.willRetain = willRetain
}

//WriteTo writes the connect control packet to writer according to the mqtt protocol.
func (c Connect) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	n1, err := writeFixedHeader(c, writer)
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write connect packet: failed to write fixed header: %v", err)
	}

	n2, err := writeVariableHeader(c, writer)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write connect packet: failed to write variable header: %v", err)
	}

	n3, err := writePayload(c, writer)
	n += n3
	if err != nil {
		return n, fmt.Errorf("failed to write connect packet: failed to write payload: %v", err)
	}

	return n, nil
}

func writeFixedHeader(c Connect, writer io.Writer) (int64, error) {
	var n int64
	firstHeaderByte := byte(CONNECT) << 4
	n1, err := writer.Write([]byte{firstHeaderByte})
	n += int64(n1)
	if err != nil {
		return n, fmt.Errorf("failed to write first header byte: %v", err)
	}

	// Remaining length
	var remainingLength uint32 = 7 + 1 + 2 // protocol name & version, flags, keep alive
	remainingLength += c.props.size()
	remainingLength += types.StringSize(c.payload.clientID)
	if c.payload.willTopic != "" {
		remainingLength += c.payload.willProps.size()
		remainingLength += types.StringSize(c.payload.willTopic)
		remainingLength += types.BinarySize(c.payload.willPayload)
	}
	if c.payload.username != "" {
		remainingLength += types.StringSize(c.payload.username)
	}
	if c.payload.password != nil {
		remainingLength += types.BinarySize(c.payload.password)
	}
	n2, err := types.WriteVarIntTo(writer, remainingLength)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write remaining packet length: %v", err)
	}

	return n, nil
}

func writeVariableHeader(c Connect, writer io.Writer) (int64, error) {
	var n int64
	// 3.1.2.1 Protocol Name
	// 3.1.2.2 Protocol Version
	n1, err := writer.Write([]byte{0, 4, 'M', 'Q', 'T', 'T', 5})
	n += int64(n1)
	if err != nil {
		return n, fmt.Errorf("failed to write protocol info: %v", err)
	}

	// 3.1.2.3 Connect Flags
	var flags byte
	// reserved
	flags |= 0
	if c.cleanStart {
		flags |= 1 << 1
	}
	if c.payload.willTopic != "" {
		flags |= 2 << 1
	}
	switch c.payload.willQoS {
	case Qos0:
	case Qos1:
		flags |= 1 << 3
	case Qos2:
		flags |= 2 << 3
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

	n2, err := writer.Write([]byte{flags})
	n += int64(n2)
	if err != nil {
		return n, fmt.Errorf("failed to write flags: %v", err)
	}

	// 3.1.2.10 Keep Alive
	n3, err := types.WriteUInt16To(writer, c.keepAlive)
	n += n3
	if err != nil {
		return n, fmt.Errorf("failed to write keep alive: %v", err)
	}

	// 3.1.2.11 Properties
	n4, err := c.props.WriteTo(writer)
	n += n4
	if err != nil {
		return n, fmt.Errorf("failed to write properties: %v", err)
	}

	return n, nil
}

func writePayload(c Connect, writer io.Writer) (int64, error) {
	var n int64
	// 3.1.3.1 ClientID
	n1, err := types.WriteStringTo(writer, c.payload.clientID)
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write connect packet: failed to write client id: %v", err)
	}

	// 3.1.3.2 Will properties
	if c.payload.willTopic != "" {
		// 3.1.3.2.1 Property length
		n2, err := c.payload.willProps.WriteTo(writer)
		n += n2
		if err != nil {
			return n, fmt.Errorf("failed to write connect packet: failed to write will properties: %v", err)
		}

		// 3.1.3.3 Will topic
		n3, err := types.WriteStringTo(writer, c.payload.willTopic)
		n += n3
		if err != nil {
			return n, fmt.Errorf("failed to write connect packet: failed to write will topic: %v", err)
		}

		// 3.1.3.4 Will payload
		n4, err := types.WriteBinaryTo(writer, c.payload.willPayload)
		n += n4
		if err != nil {
			return n, fmt.Errorf("failed to write connect packet: failed to write will payload: %v", err)
		}
	}

	// 3.1.3.5 User name
	if c.payload.username != "" {
		n5, err := types.WriteStringTo(writer, c.payload.username)
		n += n5
		if err != nil {
			return n, fmt.Errorf("failed to write connect packet: failed to write username: %v", err)
		}
	}

	// 3.1.3.6 Password
	if c.payload.password != nil {
		n6, err := types.WriteBinaryTo(writer, c.payload.password)
		n += n6
		if err != nil {
			return n, fmt.Errorf("failed to write Connect packet: failed to write password: %v", err)
		}
	}

	return n, nil
}

func readConnect(reader io.Reader, connect *Connect) error {
	// 3.1.2 Variable header
	var buf [8]byte
	if _, err := io.ReadFull(reader, buf[:]); err != nil {
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
	reserved := flags&byte(1) > 0
	cleanStart := flags&byte(1<<1) > 0
	hasWill := flags&byte(1<<2) > 0
	willQoS := flags & 24 >> 3
	willRetain := flags&byte(1<<5) > 0
	hasPassword := flags&byte(1<<6) > 0
	hasUsername := flags&byte(1<<7) > 0

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
	props, err := readProperties(reader)
	if err != nil {
		return fmt.Errorf("failed to read connect packet: failed to read will properties: %v", err)
	}
	connect.props = props

	// 3.1.3 Payload
	payload := ConnectPayload{}
	payload.SetWillRetain(willRetain)
	// 3.1.3.1 ClientID
	clientID, err := types.ReadString(reader)
	if err != nil {
		return fmt.Errorf("failed to read clientID: %v", err)
	}
	clientIDLength := len(clientID)
	if clientIDLength > 23 {
		return fmt.Errorf("malformed packet: ClientID too long (%d), must be between 1 and 23 char long", clientIDLength)
	}
	//TODO: check characters are only in a-zA-Z0-9
	payload.SetClientID(clientID)

	if hasWill {
		payload.SetWillQoS(willQoS)

		// 3.1.3.2 Will properties
		willProps, err := readProperties(reader)
		if err != nil {
			return fmt.Errorf("failed to read connect packet: failed to read will properties: %v", err)
		}
		payload.willProps = willProps

		// 3.1.3.3 Will topic
		willTopic, err := types.ReadString(reader)
		if err != nil {
			return fmt.Errorf("failed to read will topic: %v", err)
		}
		payload.SetWillTopic(willTopic)

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
