package packet

import (
	"io"
)

type Type int8

type QoS int8

const (
	Qos0 = QoS(0)
	Qos1 = QoS(1)
	Qos2 = QoS(2)
)

const (
	CONNECT Type = iota + 1
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

type Packet interface {
	Write(io.Writer) error
}
