package packet

import (
	"fmt"
	"io"
)

type pktType byte

//QoS is an alias for the three quality of service levels.
type QoS byte

//The three quality of service levels.
const (
	//QoS 0 = At most once delivery
	Qos0 = QoS(0)
	//QoS 1 = At least once delivery
	Qos1 = QoS(1)
	//QoS 2 = Exactly once delivery
	Qos2 = QoS(2)
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
//If the packet is malformec or contains a protocol error, an error is returned.
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
	switch header.pktType {
	case CONNECT:
		if header.flags != 0 {
			return nil, fmt.Errorf("failed to read packet: invalid fixed header of Connect packet: invalid flags '%d'", header.flags)
		}
		var connect Connect
		limitReader := io.LimitReader(reader, int64(header.length))
		err := readConnect(limitReader, &connect)
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
		err := readConnack(reader, &connack, header)
		if err != nil {
			return nil, fmt.Errorf("failed to read connack packet: %v", err)
		}
		return &connack, nil

	case PUBACK:
	case PUBREC:
	case PUBREL:
	case PUBCOMP:
	case SUBSCRIBE:
	case SUBACK:
	case UNSUBSCRIBE:
	case UNSUBACK:
	case PINGREQ:
	case PINGRESP:
		panic("implement me")

	case DISCONNECT:
		//if header.Flags() != 0 {
		//return nil, fmt.Errorf("failed to read packet: invalid fixed header of Disconnect packet: invalid flags '%d'", header.Flags())
		//}
		//var disconnect Disconnect
		//err := ReadDisconnect(reader, &disconnect, header)
		//if err != nil {
		//return nil, fmt.Errorf("failed to read Disconnect packet: %v", err)
		//}
		//log.Info("read Disconnect packet")
		//return &disconnect, nil

	case AUTH:
		panic("implement me")
	}
	return nil, fmt.Errorf("header with invalid packet type '%v'", header.pktType)
}
