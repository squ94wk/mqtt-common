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
			return nil, fmt.Errorf("failed to read packet: invalid fixed header of Connect packet: invalid flags '%d'", header.flags)
		}
		var connect Connect
		err := readConnect(limitedReader, &connect)
		if err != nil {
			return nil, fmt.Errorf("failed to read Connect packet: %v", err)
		}
		return &connect, nil

	case PUBLISH:
		//var publish Publish
		//err := ReadPublish(reader, &publish, header)
		//if err != nil {
		//return nil, fmt.Errorf("failed to read Publish packet: %v", err)
		//}
		//log.Info("read Publish packet")
		//return &publish, nil

	case CONNACK:
		if header.flags != 0 {
			return nil, fmt.Errorf("failed to read packet: invalid fixed header of Connack packet: invalid flags '%d'", header.flags)
		}
		var connack Connack
		err := readConnack(limitedReader, &connack)
		if err != nil {
			return nil, fmt.Errorf("failed to read connack packet: %v", err)
		}
		return &connack, nil

	case PUBACK:
	case PUBREC:
	case PUBREL:
	case PUBCOMP:
	case SUBSCRIBE:
		if header.flags != 2 {
			return nil, fmt.Errorf("failed to read packet: invalid fixed header of Subscribe packet: invalid flags '%d'", header.flags)
		}
		var subscribe Subscribe
		err := readSubscribe(limitedReader, &subscribe)
		if err != nil {
			return nil, fmt.Errorf("failed to read subscribe packet: %v", err)
		}
		return &subscribe, nil

	case SUBACK:
		if header.flags != 0 {
			return nil, fmt.Errorf("failed to read packet: invalid fixed header of Suback packet: invalid flags '%d'", header.flags)
		}
		var suback Suback
		err := readSuback(limitedReader, &suback)
		if err != nil {
			return nil, fmt.Errorf("failed to read suback packet: %v", err)
		}
		return &suback, nil

	case UNSUBSCRIBE:
	case UNSUBACK:
	case PINGREQ:
	case PINGRESP:
		panic("implement me")

	case DISCONNECT:
		if header.flags != 0 {
			return nil, fmt.Errorf("failed to read packet: invalid fixed header of Disconnect packet: invalid flags '%d'", header.flags)
		}

		var disconnect Disconnect
		err := readDisconnect(limitedReader, &disconnect, header.length)
		if err != nil {
			return nil, fmt.Errorf("failed to read disconnect packet: %v", err)
		}
		return &disconnect, nil

	case AUTH:
		panic("implement me")
	}
	return nil, fmt.Errorf("header with invalid packet type '%v'", header.pktType)
}
