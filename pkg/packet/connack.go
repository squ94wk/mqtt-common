package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

//Connack defines the connack control packet.
type Connack struct {
	sessionPresent bool
	connectReason  ConnectReason
	props          map[PropID][]Property
}

//NewConnack is the constructor for the Connack type.
func NewConnack(sessionPresent bool, connectReason ConnectReason, props map[PropID][]Property) Connack {
	return Connack{
		sessionPresent: sessionPresent,
		connectReason:  connectReason,
		props:          props,
	}
}

//SessionPresent returns the value of the session present flag.
func (c Connack) SessionPresent() bool {
	return c.sessionPresent
}

//SetSessionPresent sets the value of the session present flag.
func (c *Connack) SetSessionPresent(present bool) {
	c.sessionPresent = present
}

//ConnectReason returns the value of the connect reason.
func (c Connack) ConnectReason() ConnectReason {
	return c.connectReason
}

//SetConnectReason sets the value of the connect reason.
func (c *Connack) SetConnectReason(reason ConnectReason) {
	c.connectReason = reason
}

//Props returns the properties of the connack control packet.
func (c Connack) Props() map[PropID][]Property {
	return c.props
}

//SetProps replaces the properties of the connack control packet.
func (c *Connack) SetProps(props map[PropID][]Property) {
	c.props = props
}

//AddProp adds a property to the properties of the connack control packet.
//If the packet already contains any properties with the same property identifier it is appended the the existing ones.
//It also makes no assumptions as to if the mqtt protocol allows multiple properties of that identifier.
func (c *Connack) AddProp(prop Property) {
	propID := prop.PropID()
	properties, ok := c.props[propID]
	if !ok {
		c.props[propID] = []Property{prop}
	} else {
		c.props[propID] = append(properties, prop)
	}
}

//ResetProps removes all properties from the connack control packet
func (c *Connack) ResetProps() {
	c.props = make(map[PropID][]Property)
}

//WriteTo writes the connack control packet to writer according to the mqtt protocol.
func (c Connack) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	// 3.2.1 Fixed header
	firstByte := byte(CONNACK) << 4
	n1, err := writer.Write([]byte{firstByte})
	n += int64(n1)
	if err != nil {
		return n, fmt.Errorf("failed to write connack packet: failed to write fixed header: %v", err)
	}

	//3.2.2 Variable header
	var remainingLength uint32 = 1 + 1         // flags = session present, connect reason
	remainingLength += propertiesSize(c.props) // size of varInt of the props length
	// no payload
	n2, err := types.WriteVarIntTo(writer, uint32(remainingLength))
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write connack packet: failed to write packet length: %v", err)
	}

	var encFlags [1]byte
	if c.sessionPresent {
		encFlags[0] = 1
	} else {
		encFlags[0] = 0
	}
	n3, err := writer.Write(encFlags[:])
	n += int64(n3)
	if err != nil {
		return n, fmt.Errorf("failed to write connack packet: failed to write flags: %v", err)
	}

	connectReasonBuf := []byte{byte(c.connectReason)}
	n4, err := writer.Write(connectReasonBuf)
	n += int64(n4)
	if err != nil {
		return n, fmt.Errorf("failed to write connack packet: failed to write connect reason: %v", err)
	}

	if _, err := WritePropsTo(writer, c.props); err != nil {
		return n, fmt.Errorf("failed to write connack packet: failed to write properties: %v", err)
	}

	return n, nil
}

func readConnack(reader io.Reader, connack *Connack, header header) error {
	limitReader := io.LimitReader(reader, int64(header.length))

	// 3.2.2 Variable header
	var buf [2]byte
	_, err := io.ReadFull(limitReader, buf[:])
	if err != nil {
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
	err = readProps(limitReader, connack.props)
	if err != nil {
		return fmt.Errorf("failed to read connack packet: failed to read properties: %v", err)
	}

	return nil
}
