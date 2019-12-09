package packet

import (
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

//PropID is an alias for all defined identifiers a property can have.
type PropID uint32

//Property defines a property.
type Property interface {
	PropID() PropID
	WriteTo(writer io.Writer) (int64, error)
	size() uint32
}

type property struct {
	propID PropID
}

//ByteProp defines a byte property.
type ByteProp struct {
	property
	value byte
}

//Int32Prop defines an 32 bit integer property.
type Int32Prop struct {
	property
	value uint32
}

//Int16Prop defines an 16 bit integer property.
type Int16Prop struct {
	property
	value uint16
}

//StringProp defines a string property.
type StringProp struct {
	property
	value string
}

//KeyValueProp defines a key value property made up of two strings.
type KeyValueProp struct {
	property
	key   string
	value string
}

//VarIntProp defines a variable length integer property.
type VarIntProp struct {
	property
	value uint32
}

//BinaryProp defines a binary property.
type BinaryProp struct {
	property
	value []byte
}

//Constants for all defined property identifiers
const (
	PayloadFormatIndicator          PropID = 1
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

//PropID returns the property identifier of the property.
func (p property) PropID() PropID {
	return p.propID
}

//NewByteProp is the constructor for a byte property.
func NewByteProp(id PropID, val byte) ByteProp {
	return ByteProp{property: property{id}, value: val}
}

//NewInt32Prop is the constructor for a 32 bit integer property.
func NewInt32Prop(id PropID, val uint32) Int32Prop {
	return Int32Prop{property: property{id}, value: val}
}

//NewInt16Prop is the constructor for a 16 bit integer property.
func NewInt16Prop(id PropID, val uint16) Int16Prop {
	return Int16Prop{property: property{id}, value: val}
}

//NewStringProp is the constructor for a string property.
func NewStringProp(id PropID, val string) StringProp {
	return StringProp{property: property{id}, value: val}
}

//NewKeyValueProp is the constructor for a key value property.
func NewKeyValueProp(id PropID, key string, value string) KeyValueProp {
	return KeyValueProp{property: property{id}, key: key, value: value}
}

//NewVarIntProp is the constructor for a variable length integer property.
func NewVarIntProp(id PropID, val uint32) VarIntProp {
	return VarIntProp{property: property{id}, value: val}
}

//NewBinaryProp is the constructor for a binary property.
func NewBinaryProp(id PropID, val []byte) BinaryProp {
	return BinaryProp{property: property{id}, value: val}
}

//Value returns the value of the byte property.
func (p ByteProp) Value() byte {
	return p.value
}

//Value returns the value of the 16 bit integer property.
func (p Int16Prop) Value() uint16 {
	return p.value
}

//Value returns the value of the 32 bit integer property.
func (p Int32Prop) Value() uint32 {
	return p.value
}

//Value returns the value of the variable length integer property.
func (p VarIntProp) Value() uint32 {
	return p.value
}

//Value returns the value of the string property.
func (p StringProp) Value() string {
	return p.value
}

//Key returns the key of the key value property.
func (p KeyValueProp) Key() string {
	return p.key
}

//Value returns the value of the key value property.
func (p KeyValueProp) Value() string {
	return p.value
}

//Value returns the value of the binary property.
func (p BinaryProp) Value() []byte {
	return p.value
}

//BuildProps is an auxiliary function to collect properties together into a map.
func BuildProps(props ...Property) map[PropID][]Property {
	properties := make(map[PropID][]Property)
	for _, p := range props {
		if withID, ok := properties[p.PropID()]; ok {
			properties[p.PropID()] = append(withID, p)
		} else {
			properties[p.PropID()] = []Property{p}
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
	return 1 + uint32(2+len(p.value))
}

func (p KeyValueProp) size() uint32 {
	return 1 + uint32(2+len(p.key)+2+len(p.value))
}

func (p VarIntProp) size() uint32 {
	return 1 + types.VarIntSize(p.value)
}

func (p BinaryProp) size() uint32 {
	return 1 + uint32(len(p.value)+len(p.value))
}

func propertiesSize(props map[PropID][]Property) uint32 {
	var propLength uint32
	for _, propsForID := range props {
		for _, prop := range propsForID {
			propLength += prop.size()
		}
	}

	return types.VarIntSize(propLength) + propLength
}
