package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

//Disconnect defines the disconnect control packet.
type Disconnect struct {
	Reason DisconnectReason
	Props  Properties
}

//WriteTo writes the disconnect control packet to writer according to the mqtt protocol.
func (d Disconnect) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	// 3.14.1 Fixed header
	firstByte := byte(DISCONNECT) << 4
	n1, err := writer.Write([]byte{firstByte})
	n += int64(n1)
	if err != nil {
		return n, fmt.Errorf("failed to write disconnect packet: failed to write fixed header: %v", err)
	}

	//3.14.2 Variable header
	if d.Reason == DisconnectNormalDisconnection && len(d.Props) == 0 {
		n2, err := types.WriteVarIntTo(writer, 0)
		n += n2
		if err != nil {
			return n, fmt.Errorf("failed to write disconnect packet: failed to write packet length: %v", err)
		}
		return n, nil
	}

	if len(d.Props) == 0 {
		n2, err := types.WriteVarIntTo(writer, 1)
		n += n2
		if err != nil {
			return n, fmt.Errorf("failed to write disconnect packet: failed to write packet length: %v", err)
		}
		disconnectReasonBuf := []byte{byte(d.Reason)}
		n3, err := writer.Write(disconnectReasonBuf)
		n += int64(n3)
		if err != nil {
			return n, fmt.Errorf("failed to write disconnect packet: failed to write disconnect reason: %v", err)
		}
		return n, nil
	}

	var remainingLength uint32 = 1 // disconnect reason
	remainingLength += d.Props.size()
	// no payload
	n2, err := types.WriteVarIntTo(writer, remainingLength)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write disconnect packet: failed to write packet length: %v", err)
	}

	disconnectReasonBuf := []byte{byte(d.Reason)}
	n4, err := writer.Write(disconnectReasonBuf)
	n += int64(n4)
	if err != nil {
		return n, fmt.Errorf("failed to write disconnect packet: failed to write disconnect reason: %v", err)
	}

	n5, err := d.Props.WriteTo(writer)
	n += n5
	if err != nil {
		return n, fmt.Errorf("failed to write disconnect packet: failed to write properties: %v", err)
	}

	return n, nil
}

func readDisconnect(reader io.Reader, disconnect *Disconnect, remainingLength uint32) error {
	// 3.14.2 Variable header
	//default reason is inferred if length is 0
	if remainingLength < 1 {
		disconnect.Reason = DisconnectNormalDisconnection
		disconnect.Props = NewProperties()
		return nil
	}

	// 3.14.2.1 Disconnect reason code
	var buf [1]byte
	_, err := io.ReadFull(reader, buf[:])
	if err != nil {
		return fmt.Errorf("failed to read disconnect packet: failed to read variable header: %v", err)
	}

	//TODO: check for allowed values
	disconnect.Reason = DisconnectReason(buf[0])
	if remainingLength < 2 {
		disconnect.Props = NewProperties()
		return nil
	}

	// 3.14.2.2 Disconnect properties
	props, err := readProperties(reader)
	if err != nil {
		return fmt.Errorf("failed to read disconnect packet: failed to read properties: %v", err)
	}
	disconnect.Props = props

	return nil
}
