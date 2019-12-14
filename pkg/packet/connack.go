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
	props          Properties
}

//NewConnack is the constructor for the Connack type.
func NewConnack(sessionPresent bool, connectReason ConnectReason, props Properties) Connack {
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
func (c Connack) Props() Properties {
	return c.props
}

//SetProps replaces the properties of the connack control packet.
func (c *Connack) SetProps(props map[uint32][]Property) {
	c.props = props
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
	remainingLength += c.props.size()
	// no payload
	n2, err := types.WriteVarIntTo(writer, remainingLength)
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

	n5, err := c.props.WriteTo(writer)
	n += n5
	if err != nil {
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
	props, err := readProperties(limitReader)
	if err != nil {
		return fmt.Errorf("failed to read connack packet: failed to read properties: %v", err)
	}
	connack.props = props

	return nil
}
