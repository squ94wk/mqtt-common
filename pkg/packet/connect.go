package packet

import (
	"bytes"
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

//Connect defines the connect control packet.
type Connect struct {
	KeepAlive  uint16
	CleanStart bool
	Props      Properties
	Payload    ConnectPayload
}

//ConnectPayload defines the payload of a connect control packet.
type ConnectPayload struct {
	ClientID    string
	WillProps   Properties
	WillRetain  bool
	WillQoS     byte
	WillTopic   string
	WillPayload []byte
	Username    string
	Password    []byte
}

//WriteTo writes the connect control packet to writer according to the mqtt protocol.
func (c Connect) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	n1, err := writeFixedConnectHeader(c, writer)
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write connect packet: failed to write fixed header: %v", err)
	}

	n2, err := writeVariableConnectHeader(c, writer)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write connect packet: failed to write variable header: %v", err)
	}

	n3, err := writeConnectPayload(c, writer)
	n += n3
	if err != nil {
		return n, fmt.Errorf("failed to write connect packet: failed to write payload: %v", err)
	}

	return n, nil
}

func writeFixedConnectHeader(c Connect, writer io.Writer) (int64, error) {
	var n int64
	firstHeaderByte := byte(CONNECT) << 4
	n1, err := writer.Write([]byte{firstHeaderByte})
	n += int64(n1)
	if err != nil {
		return n, fmt.Errorf("failed to write first header byte: %v", err)
	}

	// Remaining length
	var remainingLength uint32 = 7 + 1 + 2 // protocol name & version, flags, keep alive
	remainingLength += c.Props.size()
	remainingLength += types.StringSize(c.Payload.ClientID)
	if c.Payload.WillTopic != "" {
		remainingLength += c.Payload.WillProps.size()
		remainingLength += types.StringSize(c.Payload.WillTopic)
		remainingLength += types.BinarySize(c.Payload.WillPayload)
	}
	if c.Payload.Username != "" {
		remainingLength += types.StringSize(c.Payload.Username)
	}
	if c.Payload.Password != nil {
		remainingLength += types.BinarySize(c.Payload.Password)
	}
	n2, err := types.WriteVarIntTo(writer, remainingLength)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write remaining packet length: %v", err)
	}

	return n, nil
}

func writeVariableConnectHeader(c Connect, writer io.Writer) (int64, error) {
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
	if c.CleanStart {
		flags |= 1 << 1
	}
	if c.Payload.WillTopic != "" {
		flags |= 2 << 1
	}
	switch c.Payload.WillQoS {
	case Qos0:
	case Qos1:
		flags |= 1 << 3
	case Qos2:
		flags |= 2 << 3
	default:
		panic(fmt.Errorf("invalid option for QoS: %v", c.Payload.WillQoS))
	}
	if c.Payload.WillRetain {
		flags |= 1 << 5
	}
	if c.Payload.Password != nil {
		flags |= 1 << 6
	}
	if c.Payload.Username != "" {
		flags |= 1 << 7
	}

	n2, err := writer.Write([]byte{flags})
	n += int64(n2)
	if err != nil {
		return n, fmt.Errorf("failed to write flags: %v", err)
	}

	// 3.1.2.10 Keep Alive
	n3, err := types.WriteUInt16To(writer, c.KeepAlive)
	n += n3
	if err != nil {
		return n, fmt.Errorf("failed to write keep alive: %v", err)
	}

	// 3.1.2.11 Properties
	n4, err := c.Props.WriteTo(writer)
	n += n4
	if err != nil {
		return n, fmt.Errorf("failed to write properties: %v", err)
	}

	return n, nil
}

func writeConnectPayload(c Connect, writer io.Writer) (int64, error) {
	var n int64
	// 3.1.3.1 ClientID
	n1, err := types.WriteStringTo(writer, c.Payload.ClientID)
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write connect packet: failed to write client id: %v", err)
	}

	// 3.1.3.2 Will properties
	if c.Payload.WillTopic != "" {
		// 3.1.3.2.1 Property length
		n2, err := c.Payload.WillProps.WriteTo(writer)
		n += n2
		if err != nil {
			return n, fmt.Errorf("failed to write connect packet: failed to write will properties: %v", err)
		}

		// 3.1.3.3 Will topic
		n3, err := types.WriteStringTo(writer, c.Payload.WillTopic)
		n += n3
		if err != nil {
			return n, fmt.Errorf("failed to write connect packet: failed to write will topic: %v", err)
		}

		// 3.1.3.4 Will payload
		n4, err := types.WriteBinaryTo(writer, c.Payload.WillPayload)
		n += n4
		if err != nil {
			return n, fmt.Errorf("failed to write connect packet: failed to write will payload: %v", err)
		}
	}

	// 3.1.3.5 User name
	if c.Payload.Username != "" {
		n5, err := types.WriteStringTo(writer, c.Payload.Username)
		n += n5
		if err != nil {
			return n, fmt.Errorf("failed to write connect packet: failed to write username: %v", err)
		}
	}

	// 3.1.3.6 Password
	if c.Payload.Password != nil {
		n6, err := types.WriteBinaryTo(writer, c.Payload.Password)
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

	connect.CleanStart = cleanStart

	// 3.1.2.6 Will QoS
	if willQoS > 2 {
		return fmt.Errorf("malformed packet: Will QoS must be 0, 1 or 2, but is %d", willQoS)
	}

	// 3.1.2.10 Keep Alive
	keepAlive, err := types.ReadUInt16(reader)
	if err != nil {
		return fmt.Errorf("failed to read connect packet: failed to read keepAlive: %v", err)
	}
	connect.KeepAlive = keepAlive

	// 3.1.2.11 Properties
	props, err := readProperties(reader)
	if err != nil {
		return fmt.Errorf("failed to read connect packet: failed to read properties: %v", err)
	}
	connect.Props = props

	// 3.1.3 Payload
	payload := ConnectPayload{}
	payload.WillRetain = willRetain
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
	payload.ClientID = clientID

	if hasWill {
		payload.WillQoS = willQoS

		// 3.1.3.2 Will properties
		willProps, err := readProperties(reader)
		if err != nil {
			return fmt.Errorf("failed to read connect packet: failed to read will properties: %v", err)
		}
		payload.WillProps = willProps

		// 3.1.3.3 Will topic
		willTopic, err := types.ReadString(reader)
		if err != nil {
			return fmt.Errorf("failed to read will topic: %v", err)
		}
		payload.WillTopic = willTopic

		// 3.1.3.4 Will payload
		willPayload, err := types.ReadBinary(reader)
		if err != nil {
			return fmt.Errorf("failed to read will payload: %v", err)
		}
		payload.WillPayload = willPayload
	}

	// 3.1.3.5 User name
	if hasUsername {
		username, err := types.ReadString(reader)
		if err != nil {
			return fmt.Errorf("failed to read username: %v", err)
		}
		payload.Username = username
	}

	// 3.1.3.6 Password
	if hasPassword {
		password, err := types.ReadBinary(reader)
		if err != nil {
			return fmt.Errorf("failed to read password: %v", err)
		}
		payload.Password = password
	}

	connect.Payload = payload
	return nil
}
