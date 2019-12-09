package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

type propReader func(io.Reader, PropID) (Property, error)

var (
	propReaders map[PropID]propReader
)

func init() {
	propReaders = make(map[PropID]propReader)

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

func readProp(reader io.Reader) (Property, error) {
	val, err := types.ReadVarInt(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read property: failed to read property identifier: '%v'", err)
	}

	propID := PropID(val)
	propReader, ok := propReaders[propID]
	if !ok {
		//TODO: panic, if we place check beforehand?
		return nil, fmt.Errorf("failed to read property: no reader for property with identifier '%d'", propID)
	}

	property, err := propReader(reader, propID)
	if err != nil {
		return nil, fmt.Errorf("failed to read property with identifier '%d': %v", propID, err)
	}

	return property, nil
}

func readProps(reader io.Reader, props map[PropID][]Property) error {
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
		property, err := readProp(limitReader)
		if err != nil {
			return fmt.Errorf("failed to read property: %v", err)
		}

		properties, ok := props[property.PropID()]
		if ok {
			props[property.PropID()] = append(properties, property)
		} else {
			props[property.PropID()] = []Property{property}
		}
	}

	return nil
}

func readByteProp(reader io.Reader, propID PropID) (Property, error) {
	var buf [1]byte
	if _, err := io.ReadFull(reader, buf[:1]); err != nil {
		return nil, fmt.Errorf("failed to read byte property: %v", err)
	}
	return ByteProp{value: buf[0], property: property{propID}}, nil
}

func readInt32Prop(reader io.Reader, propID PropID) (Property, error) {
	val, err := types.ReadUInt32(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read four byte integer property: %v", err)
	}

	return Int32Prop{value: val, property: property{propID}}, nil
}

func readInt16Prop(reader io.Reader, propID PropID) (Property, error) {
	val, err := types.ReadUInt16(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read two byte integer property: %v", err)
	}

	return Int16Prop{value: val, property: property{propID}}, nil
}

func readStringProp(reader io.Reader, propID PropID) (Property, error) {
	val, err := types.ReadString(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read string property: %v", err)
	}

	return StringProp{value: val, property: property{propID}}, nil
}

func readKeyValueProp(reader io.Reader, propID PropID) (Property, error) {
	key, err := types.ReadString(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read string pair property: failed to read key: %v", err)
	}
	value, err := types.ReadString(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read string pair property: failed to read value: %v", err)
	}

	return KeyValueProp{key: key, value: value, property: property{propID}}, nil
}

func readVarIntProp(reader io.Reader, propID PropID) (Property, error) {
	val, err := types.ReadVarInt(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read variable length int property: %v", err)
	}

	return VarIntProp{value: val, property: property{propID}}, nil
}

func readBinaryProp(reader io.Reader, propID PropID) (Property, error) {
	val, err := types.ReadBinary(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read string pair property: %v", err)
	}

	return BinaryProp{value: val, property: property{propID}}, nil
}
