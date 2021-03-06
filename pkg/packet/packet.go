package packet

import (
	"fmt"
	"io"
)

type pktType byte

//The three quality of service levels.
const (
	//QoS 0 = At most once delivery
	Qos0          = 0
	QosAtMostOnce = 0
	//QoS 1 = At least once delivery
	Qos1           = 1
	QosAtLeastOnce = 1
	//QoS 2 = Exactly once delivery
	Qos2           = 2
	QosExactlyOnce = 2
)

//Constants for all control packet types
const (
	CONNECT pktType = iota + 1
	CONNACK
	PUBLISH
	PUBACK
	PUBREC
	PUBREL
	PUBCOMP
	SUBSCRIBE
	SUBACK
	UNSUBSCRIBE
	UNSUBACK
	PINGREQ
	PINGRESP
	DISCONNECT
	AUTH
)

//Packet defines a control packet.
type Packet interface {
	WriteTo(io.Writer) (int64, error)
}

//ReadPacket reads a packet from reader.
//If the packet is malformed or contains a protocol error, an error is returned.
func ReadPacket(reader io.Reader) (Packet, error) {
	var header header
	if err := readHeader(reader, &header); err != nil {
		return nil, fmt.Errorf("failed to read packet: failed to read header: %v", err)
	}

	pkt, err := readRestOfPacket(reader, header)
	if err != nil {
		return nil, err
	}

	return pkt, nil
}

func readRestOfPacket(reader io.Reader, header header) (Packet, error) {
	limitedReader := io.LimitReader(reader, int64(header.length))
	switch header.pktType {
	case CONNECT:
		if header.flags != 0 {
			return nil, fmt.Errorf("failed to read Connect packet: invalid fixed header: invalid flags '%d'", header.flags)
		}
		var connect Connect
		err := readConnect(limitedReader, &connect)
		if err != nil {
			return nil, fmt.Errorf("failed to read Connect packet: %v", err)
		}
		return &connect, nil

	case PUBLISH:
		var publish Publish
		err := readPublish(limitedReader.(*io.LimitedReader), &publish, header.flags)
		if err != nil {
			return nil, fmt.Errorf("failed to read Publish packet: %v", err)
		}
		return &publish, nil

	case CONNACK:
		if header.flags != 0 {
			return nil, fmt.Errorf("failed to read Connack packet: invalid fixed header: invalid flags '%d'", header.flags)
		}
		var connack Connack
		err := readConnack(limitedReader, &connack)
		if err != nil {
			return nil, fmt.Errorf("failed to read Connack packet: %v", err)
		}
		return &connack, nil

	case SUBSCRIBE:
		if header.flags != 2 {
			return nil, fmt.Errorf("failed to read Sbuscribe packet: invalid fixed header: invalid flags '%d'", header.flags)
		}
		var subscribe Subscribe
		err := readSubscribe(limitedReader, &subscribe)
		if err != nil {
			return nil, fmt.Errorf("failed to read Subscribe packet: %v", err)
		}
		return &subscribe, nil

	case SUBACK:
		if header.flags != 0 {
			return nil, fmt.Errorf("failed to read Suback packet: invalid fixed header: invalid flags '%d'", header.flags)
		}
		var suback Suback
		err := readSuback(limitedReader, &suback)
		if err != nil {
			return nil, fmt.Errorf("failed to read Suback packet: %v", err)
		}
		return &suback, nil

	case PUBACK:
		fallthrough
	case PUBREC:
		fallthrough
	case PUBREL:
		fallthrough
	case PUBCOMP:
		fallthrough
	case UNSUBSCRIBE:
		fallthrough
	case UNSUBACK:
		fallthrough
	case PINGREQ:
		fallthrough
	case PINGRESP:
		panic("implement me")

	case DISCONNECT:
		if header.flags != 0 {
			return nil, fmt.Errorf("failed to read Disconnect packet: invalid fixed header: invalid flags '%d'", header.flags)
		}
		var disconnect Disconnect
		err := readDisconnect(limitedReader, &disconnect, header.length)
		if err != nil {
			return nil, fmt.Errorf("failed to read Disconnect packet: %v", err)
		}
		return &disconnect, nil

	case AUTH:
		panic("implement me")
	}
	return nil, fmt.Errorf("header with invalid packet type '%v'", header.pktType)
}
