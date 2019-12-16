package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/pkg/types"
)

type propReader func(io.Reader, PropId) (Property, error)

var (
	propReaders map[PropId]propReader
)

func init() {
	propReaders = make(map[PropId]propReader)

	propReaders[PayloadFormatIndicator] = readByteProp
	propReaders[MessageExpiryInterval] = readInt32Prop
	propReaders[ContentType] = readStringProp
	propReaders[ResponseTopic] = readStringProp
	propReaders[CorrelationData] = readBinaryProp
	propReaders[SubscriptionIdentifier] = readVarIntProp
	propReaders[SessionExpiryInterval] = readInt32Prop
	propReaders[AssignedClientIdentifier] = readStringProp
	propReaders[ServerKeepAlive] = readInt16Prop
	propReaders[AuthenticationMethod] = readStringProp
	propReaders[AuthenticationData] = readBinaryProp
	propReaders[RequestProblemInformation] = readByteProp
	propReaders[WillDelayInterval] = readInt32Prop
	propReaders[RequestResponseInformation] = readByteProp
	propReaders[ResponseInformation] = readStringProp
	propReaders[ServerReference] = readStringProp
	propReaders[ReasonString] = readStringProp
	propReaders[ReceiveMaximum] = readInt16Prop
	propReaders[TopicAliasMaximum] = readInt16Prop
	propReaders[TopicAlias] = readInt16Prop
	propReaders[MaximumQoS] = readByteProp
	propReaders[RetainAvailable] = readByteProp
	propReaders[UserProperty] = readKeyValueProp
	propReaders[MaximumPacketSize] = readInt32Prop
	propReaders[WildcardSubscriptionAvailable] = readByteProp
	propReaders[SubscriptionIdentifierAvailable] = readByteProp
	propReaders[SharedSubscriptionAvailable] = readByteProp
}

func readProperty(reader io.Reader) (Property, error) {
	val, err := types.ReadVarInt(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read property: failed to read property identifier: '%v'", err)
	}

	propId := PropId(val)
	propReader, ok := propReaders[propId]
	if !ok {
		//TODO: panic, if we place check beforehand?
		return nil, fmt.Errorf("failed to read property: no reader for property with identifier '%d'", propId)
	}

	property, err := propReader(reader, propId)
	if err != nil {
		return nil, fmt.Errorf("failed to read property with identifier '%d': %v", propId, err)
	}

	return property, nil
}

func readProperties(reader io.Reader, props map[PropId][]Property) error {
	if props == nil {
		return fmt.Errorf("props must not be nil")
	}

	propLength, err := types.ReadVarInt(reader)
	if err != nil {
		return fmt.Errorf("failed to read length: %v", err)
	}

	if propLength == 0 {
		return nil
	}

	limitReader := io.LimitReader(reader, int64(propLength)).(*io.LimitedReader)
	for limitReader.N > 0 {
		property, err := readProperty(limitReader)
		if err != nil {
			return fmt.Errorf("failed to read property: %v", err)
		}

		properties, ok := props[property.PropId()]
		if ok {
			props[property.PropId()] = append(properties, property)
		} else {
			props[property.PropId()] = []Property{property}
		}
	}

	return nil
}

func readByteProp(reader io.Reader, propId PropId) (Property, error) {
	var buf [1]byte
	if _, err := io.ReadFull(reader, buf[:1]); err != nil {
		return nil, fmt.Errorf("failed to read byte property: %v", err)
	}
	return ByteProp{payload: buf[0], property: property{propId}}, nil
}

func readInt32Prop(reader io.Reader, propId PropId) (Property, error) {
	val, err := types.ReadUInt32(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read four byte integer property: %v", err)
	}

	return Int32Prop{payload: val, property: property{propId}}, nil
}

func readInt16Prop(reader io.Reader, propId PropId) (Property, error) {
	val, err := types.ReadUInt16(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read two byte integer property: %v", err)
	}

	return Int16Prop{payload: val, property: property{propId}}, nil
}

func readStringProp(reader io.Reader, propId PropId) (Property, error) {
	val, err := types.ReadString(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read string property: %v", err)
	}

	return StringProp{payload: val, property: property{propId}}, nil
}

func readKeyValueProp(reader io.Reader, propId PropId) (Property, error) {
	val, err := types.ReadStringPair(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read string pair property: %v", err)
	}

	return KeyValueProp{payload: val, property: property{propId}}, nil
}

func readVarIntProp(reader io.Reader, propId PropId) (Property, error) {
	val, err := types.ReadVarInt(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read variable length int property: %v", err)
	}

	return VarIntProp{payload: val, property: property{propId}}, nil
}

func readBinaryProp(reader io.Reader, propId PropId) (Property, error) {
	val, err := types.ReadBinary(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read string pair property: %v", err)
	}

	return BinaryProp{payload: val, property: property{propId}}, nil
}
