package packet

import (
	"io"

	"github.com/squ94wk/mqtt-common/pkg/types"
)

type PropId uint32

type Property interface {
	PropId() PropId
	Write(writer io.Writer) error
	size() uint32
}

type property struct {
	propId PropId
}

type ByteProp struct {
	property
	payload byte
}

type Int32Prop struct {
	property
	payload uint32
}

type Int16Prop struct {
	property
	payload uint16
}

type StringProp struct {
	property
	payload string
}

type KeyValueProp struct {
	property
	payload types.StringPair
}

type VarIntProp struct {
	property
	payload uint32
}

type BinaryProp struct {
	property
	payload []byte
}

type properties struct {
	props map[PropId][]Property
}

const (
	PayloadFormatIndicator          PropId = 1
	MessageExpiryInterval                  = 2
	ContentType                            = 3
	ResponseTopic                          = 8
	CorrelationData                        = 9
	SubscriptionIdentifier                 = 11
	SessionExpiryInterval                  = 17
	AssignedClientIdentifier               = 18
	ServerKeepAlive                        = 19
	AuthenticationMethod                   = 21
	AuthenticationData                     = 22
	RequestProblemInformation              = 23
	WillDelayInterval                      = 24
	RequestResponseInformation             = 25
	ResponseInformation                    = 26
	ServerReference                        = 28
	ReasonString                           = 31
	ReceiveMaximum                         = 33
	TopicAliasMaximum                      = 34
	TopicAlias                             = 35
	MaximumQoS                             = 36
	RetainAvailable                        = 37
	UserProperty                           = 38
	MaximumPacketSize                      = 39
	WildcardSubscriptionAvailable          = 40
	SubscriptionIdentifierAvailable        = 41
	SharedSubscriptionAvailable            = 42
)

func (p property) PropId() PropId {
	return p.propId
}

func NewByteProp(id PropId, val byte) ByteProp {
	return ByteProp{property: property{id}, payload: val}
}

func NewInt32Prop(id PropId, val uint32) Int32Prop {
	return Int32Prop{property: property{id}, payload: val}
}

func NewInt16Prop(id PropId, val uint16) Int16Prop {
	return Int16Prop{property: property{id}, payload: val}
}

func NewStringProp(id PropId, val string) StringProp {
	return StringProp{property: property{id}, payload: val}
}

func NewKeyValueProp(id PropId, key string, value string) KeyValueProp {
	return KeyValueProp{property: property{id}, payload: types.NewStringPair(key, value)}
}

func NewVarIntProp(id PropId, val uint32) VarIntProp {
	return VarIntProp{property: property{id}, payload: val}
}

func NewBinaryProp(id PropId, val []byte) BinaryProp {
	return BinaryProp{property: property{id}, payload: val}
}

func BuildProps(props ...Property) map[PropId][]Property {
	properties := make(map[PropId][]Property)
	for _, p := range props {
		if withId, ok := properties[p.PropId()]; ok {
			properties[p.PropId()] = append(withId, p)
		} else {
			properties[p.PropId()] = []Property{p}
		}
	}

	return properties
}

func (p ByteProp) size() uint32 {
	return 1 + 1
}

func (p Int32Prop) size() uint32 {
	return 1 + 4
}

func (p Int16Prop) size() uint32 {
	return 1 + 2
}

func (p StringProp) size() uint32 {
	return 1 + uint32(2+len(p.payload))
}

func (p KeyValueProp) size() uint32 {
	return 1 + uint32(2+len(p.payload.Key())+2+len(p.payload.Value()))
}

func (p VarIntProp) size() uint32 {
	return 1 + types.VarIntSize(p.payload)
}

func (p BinaryProp) size() uint32 {
	return 1 + uint32(len(p.payload)+len(p.payload))
}

func propertiesSize(props map[PropId][]Property) uint32 {
	var propLength uint32 = 0
	for _, propsForId := range props {
		for _, prop := range propsForId {
			propLength += uint32(prop.size())
		}
	}

	return types.VarIntSize(propLength) + propLength
}

func (p *properties) ResetProps() {
	p.props = make(map[PropId][]Property)
}

func (p *properties) SetProps(props map[PropId][]Property) {
	p.props = props
}

func (p *properties) AddProp(prop Property) {
	propId := prop.PropId()

	properties, ok := p.props[propId]
	if !ok {
		p.props[propId] = []Property{prop}
	} else {
		p.props[propId] = append(properties, prop)
	}
}
