package packet

import (
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

//Property defines a property composition.
type Property struct {
	propID  uint32
	payload PropertyPayload
}

//PropertyPayload defines the part of a property that is specific to the type of property.
type PropertyPayload interface {
	WriteTo(writer io.Writer) (int64, error)
	size() uint32
}

//Properties is an auxiliary type that holds many properties.
type Properties map[uint32][]Property

//BytePropPayload defines a byte property.
type BytePropPayload byte

//Int32PropPayload defines an 32 bit integer property.
type Int32PropPayload uint32

//Int16PropPayload defines an 16 bit integer property.
type Int16PropPayload uint16

//StringPropPayload defines a string property.
type StringPropPayload string

//KeyValuePropPayload defines a key value property made up of two strings.
type KeyValuePropPayload struct {
	key   string
	value string
}

//VarIntPropPayload defines a variable length integer property.
type VarIntPropPayload uint32

//BinaryPropPayload defines a binary property.
type BinaryPropPayload []byte

//Constants for all defined property identifiers
const (
	PayloadFormatIndicator          uint32 = 1
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

//NewProperty constructs a new Property.
func NewProperty(propID uint32, payload PropertyPayload) Property {
	return Property{
		propID:  propID,
		payload: payload,
	}
}

//NewProperties is the constructor of the properties type.
func NewProperties(props ...Property) Properties {
	properties := make(map[uint32][]Property)
	for _, p := range props {
		if withID, ok := properties[p.PropID()]; ok {
			properties[p.PropID()] = append(withID, p)
		} else {
			properties[p.PropID()] = []Property{p}
		}
	}
	return properties
}

//NewKeyValuePropPayload is the constructor for a key value property.
func NewKeyValuePropPayload(key string, value string) KeyValuePropPayload {
	return KeyValuePropPayload{key: key, value: value}
}

//PropID returns the property identifier of the property.
func (p Property) PropID() uint32 {
	return p.propID
}

//Payload returns the payload of the property.
func (p Property) Payload() PropertyPayload {
	return p.payload
}

//Key returns the key of the key value property.
func (p KeyValuePropPayload) Key() string {
	return p.key
}

//Value returns the value of the key value property.
func (p KeyValuePropPayload) Value() string {
	return p.value
}

//Add adds a property to p.
//The property is appended to the existing ones if p already contains properties with the same identifier.
//Add makes no assumptions as to if the mqtt protocol allows multiple properties of that identifier.
func (p Properties) Add(prop Property) {
	propID := prop.PropID()
	properties, ok := p[propID]
	if !ok {
		p[propID] = []Property{prop}
	} else {
		p[propID] = append(properties, prop)
	}
}

//Reset removes all properties from p.
func (p Properties) Reset() {
	for propID := range p {
		delete(p, propID)
	}
}

func (p Property) size() uint32 {
	return types.VarIntSize(p.propID) + p.payload.size()
}

func (p BytePropPayload) size() uint32 {
	return 1
}

func (p Int32PropPayload) size() uint32 {
	return 4
}

func (p Int16PropPayload) size() uint32 {
	return 2
}

func (p StringPropPayload) size() uint32 {
	return uint32(2 + len(p))
}

func (p KeyValuePropPayload) size() uint32 {
	return uint32(2 + len(p.key) + 2 + len(p.value))
}

func (p VarIntPropPayload) size() uint32 {
	return types.VarIntSize(uint32(p))
}

func (p BinaryPropPayload) size() uint32 {
	return uint32(len(p) + len(p))
}

func (p Properties) size() uint32 {
	var propLength uint32
	for _, propsForID := range p {
		for _, prop := range propsForID {
			propLength += prop.size()
		}
	}
	return types.VarIntSize(propLength) + propLength
}
