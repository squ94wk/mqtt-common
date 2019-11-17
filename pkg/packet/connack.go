package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

type Connack struct {
	sessionPresent bool
	connectReason  ConnectReason
	props          map[PropId][]Property
}

func (c Connack) Props() map[PropId][]Property {
	return c.props
}

func (c Connack) SessionPresent() bool {
	return c.sessionPresent
}

func (c Connack) ConnectReason() ConnectReason {
	return c.connectReason
}

func (c *Connack) SetSessionPresent(present bool) {
	c.sessionPresent = present
}

func (c *Connack) SetConnectReason(reason ConnectReason) {
	c.connectReason = reason
}

func (c *Connack) SetProps(props map[PropId][]Property) {
	c.props = props
}

func (c *Connack) AddProp(prop Property) {
	propId := prop.PropId()

	properties, ok := c.props[propId]
	if !ok {
		c.props[propId] = []Property{prop}
	} else {
		c.props[propId] = append(properties, prop)
	}
}

func (c *Connack) ResetProps() {
	c.props = make(map[PropId][]Property)
}

func ReadConnack(reader io.Reader, connack *Connack, header header) error {
	limitReader := io.LimitReader(reader, int64(header.Length()))

	// 3.2.2 Variable header
	var buf [2]byte
	if _, err := io.ReadFull(limitReader, buf[:2]); err != nil {
		return fmt.Errorf("failed to read connack packet: failed to read variable header: %v", err)
	}

	// 3.2.2.1 Connect acknowledgement flags
	if buf[0] > 1 {
		return fmt.Errorf("failed to read connack packet: invalid value for flags: bits 7-1 are reserved and must be 0: got: '%v'", buf[0])
	}
	// 3.2.2.1.1 Session present
	connack.SetSessionPresent(buf[0] == 1)

	// 3.2.2.2 Connect reason code
	connack.SetConnectReason(ConnectReason(buf[1]))
	//TODO: check for allowed values

	// 3.2.2.3 Connack properties
	connack.ResetProps()
	err := readProperties(limitReader, connack.props)
	if err != nil {
		return fmt.Errorf("failed to read connack packet: failed to read properties: %v", err)
	}

	return nil
}

func (c Connack) Write(writer io.Writer) error {
	// 3.2.1 Fixed header
	firstHeaderByte := byte(CONNACK) << 4
	if _, err := writer.Write([]byte{firstHeaderByte}); err != nil {
		return fmt.Errorf("failed to write connack packet: failed to write fixed header: %v", err)
	}

	//3.2.2 Variable header
	var remainingLength uint32 = 1 + 1         // flags = session present, connect reason
	remainingLength += propertiesSize(c.props) // size of varInt of the props length
	// no payload
	if err := types.WriteVarInt(writer, uint32(remainingLength)); err != nil {
		return fmt.Errorf("failed to write connack packet: failed to write packet length: %v", err)
	}

	var encFlags [1]byte
	if c.sessionPresent {
		encFlags[0] = 1
	} else {
		encFlags[0] = 0
	}
	if _, err := writer.Write(encFlags[:1]); err != nil {
		return fmt.Errorf("failed to write connack packet: failed to write flags: %v", err)
	}

	connectReasonBuf := []byte{byte(c.connectReason)}
	if _, err := writer.Write(connectReasonBuf); err != nil {
		return fmt.Errorf("failed to write connack packet: failed to write connect reason: %v", err)
	}

	if err := WriteProperties(writer, c.props); err != nil {
		return fmt.Errorf("failed to write connack packet: failed to write properties: %v", err)
	}

	return nil
}
