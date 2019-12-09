package packet

import (
	"bytes"
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
	"github.com/squ94wk/mqtt-common/pkg/topic"
)

//Connect defines the connect control packet.
type Connect struct {
	keepAlive  uint16
	cleanStart bool
	props      map[PropID][]Property
	payload    ConnectPayload
}

//ConnectPayload defines the payload of a connect control packet.
type ConnectPayload struct {
	clientID    string
	willProps   map[PropID][]Property
	willRetain  bool
	willQoS     QoS
	willTopic   topic.Topic
	willPayload []byte
	username    string
	password    []byte
}

//NewConnect is the constructor for the Connect type.
func NewConnect(
	cleanStart bool,
	keepAlive uint16,
	props map[PropID][]Property,
	clientID string,
	willRetain bool,
	willQoS QoS,
	willProps map[PropID][]Property,
	willTopic topic.Topic,
	willPayload []byte,
	username string,
	password []byte) Connect {
	return Connect{
		keepAlive:  keepAlive,
		cleanStart: cleanStart,
		props:      props,
		payload: ConnectPayload{
			clientID:    clientID,
			willProps:   willProps,
			willRetain:  willRetain,
			willQoS:     willQoS,
			willTopic:   willTopic,
			willPayload: willPayload,
			username:    username,
			password:    password,
		},
	}
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
func (c Connect) Props() map[PropID][]Property {
	return c.props
}

//SetProps replaces the properties of the connect control packet.
func (c *Connect) SetProps(props map[PropID][]Property) {
	c.props = props
}

//AddProp adds a property to the properties of the connect control packet.
//If the packet already contains any properties with the same property identifier it is appended the the existing ones.
//It also makes no assumptions as to if the mqtt protocol allows multiple properties of that identifier.
func (c *Connect) AddProp(prop Property) {
	propID := prop.PropID()
	properties, ok := c.props[propID]
	if !ok {
		c.props[propID] = []Property{prop}
	} else {
		c.props[propID] = append(properties, prop)
	}
}

//ResetProps removes all properties from the connect control packet.
func (c *Connect) ResetProps() {
	c.props = make(map[PropID][]Property)
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
func (p ConnectPayload) WillProps() map[PropID][]Property {
	return p.willProps
}

//SetWillProps replaces the properties of the will message of the connect control packet.
func (p *ConnectPayload) SetWillProps(props map[PropID][]Property) {
	p.willProps = props
}

//AddWillProp adds a property to the properties of the connect control packet.
//If the packet already contains any properties with the same property identifier it is appended the the existing ones.
//It also makes no assumptions as to if the mqtt protocol allows multiple properties of that identifier.
func (p *ConnectPayload) AddWillProp(prop Property) {
	propID := prop.PropID()
	properties, ok := p.willProps[propID]
	if !ok {
		p.willProps[propID] = []Property{prop}
	} else {
		p.willProps[propID] = append(properties, prop)
	}
}

//ResetWillProps removes all properties from the will message of the connect control packet.
func (p *ConnectPayload) ResetWillProps() {
	p.willProps = make(map[PropID][]Property)
}

//WillTopic returns the topic of the will message of the connect control packet.
func (p ConnectPayload) WillTopic() topic.Topic {
	return p.willTopic
}

//SetWillTopic sets the topic of the will message of the connect control packet.
func (p *ConnectPayload) SetWillTopic(topic topic.Topic) {
	p.willTopic = topic
}

//WillQoS returns the quality of service level of the will message of the connect control packet.
func (p ConnectPayload) WillQoS() QoS {
	return p.willQoS
}

//SetWillQoS sets the quality of service level of the will message of the connect control packet.
func (p *ConnectPayload) SetWillQoS(qos QoS) {
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
	remainingLength += propertiesSize(c.props)
	remainingLength += types.StringSize(c.payload.clientID)
	if c.payload.willTopic != "" {
		remainingLength += propertiesSize(c.payload.willProps)
		remainingLength += types.StringSize(string(c.payload.willTopic))
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
	n4, err := WritePropsTo(writer, c.props)
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
		n2, err := WritePropsTo(writer, c.payload.willProps)
		n += n2
		if err != nil {
			return n, fmt.Errorf("failed to write connect packet: failed to write will properties: %v", err)
		}

		// 3.1.3.3 Will topic
		n3, err := types.WriteStringTo(writer, string(c.payload.willTopic))
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
	reserved := flags&byte(1) >= 1
	cleanStart := flags&byte(1<<1) >= 1
	hasWill := flags&byte(1<<2) >= 1
	willQoS := QoS(flags & 24 >> 3)
	willRetain := (flags)&byte(1<<5) >= 1
	hasPassword := flags&byte(1<<6) >= 1
	hasUsername := flags&byte(1<<7) >= 1

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
	if err := readProps(reader, connect.props); err != nil {
		return fmt.Errorf("failed to read connect packet: failed to read will properties: %v", err)
	}

	// 3.1.3 Payload
	payload := ConnectPayload{}
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

		payload.ResetWillProps()
		// 3.1.3.2 Will properties
		err := readProps(reader, payload.willProps)
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
