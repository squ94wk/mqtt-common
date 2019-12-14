package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

type propReader func(io.Reader) (PropertyPayload, error)

var (
	propReaders = map[uint32]propReader{
		PayloadFormatIndicator: readByteProp,
		MessageExpiryInterval: readInt32Prop,
		ContentType: readStringProp,
		ResponseTopic: readStringProp,
		CorrelationData: readBinaryProp,
		SubscriptionIdentifier: readVarIntProp,
		SessionExpiryInterval: readInt32Prop,
		AssignedClientIdentifier: readStringProp,
		ServerKeepAlive: readInt16Prop,
		AuthenticationMethod: readStringProp,
		AuthenticationData: readBinaryProp,
		RequestProblemInformation: readByteProp,
		WillDelayInterval: readInt32Prop,
		RequestResponseInformation: readByteProp,
		ResponseInformation: readStringProp,
		ServerReference: readStringProp,
		ReasonString: readStringProp,
		ReceiveMaximum: readInt16Prop,
		TopicAliasMaximum: readInt16Prop,
		TopicAlias: readInt16Prop,
		MaximumQoS: readByteProp,
		RetainAvailable: readByteProp,
		UserProperty: readKeyValueProp,
		MaximumPacketSize: readInt32Prop,
		WildcardSubscriptionAvailable: readByteProp,
		SubscriptionIdentifierAvailable: readByteProp,
		SharedSubscriptionAvailable: readByteProp,
	}
)

func readProp(reader io.Reader) (Property, error) {
	var prop Property
	propID, err := types.ReadVarInt(reader)
	if err != nil {
		return prop, fmt.Errorf("failed to read property: failed to read property identifier: '%v'", err)
	}
	prop.propID = propID

	propReader, ok := propReaders[propID]
	if !ok {
		//TODO: panic, if we place check beforehand?
		return prop, fmt.Errorf("failed to read property: no reader for property with identifier '%d'", propID)
	}

	payload, err := propReader(reader)
	if err != nil {
		return prop, fmt.Errorf("failed to read property with identifier '%d': %v", propID, err)
	}

	prop.payload = payload
	return prop, nil
}

func readProperties(reader io.Reader) (Properties, error) {
	props := Properties(make(map[uint32][]Property))
	propLength, err := types.ReadVarInt(reader)
	if err != nil {
		return props, fmt.Errorf("failed to read length: %v", err)
	}

	if propLength == 0 {
		return props, nil
	}

	limitReader := io.LimitReader(reader, int64(propLength)).(*io.LimitedReader)
	for limitReader.N > 0 {
		property, err := readProp(limitReader)
		if err != nil {
			return props, fmt.Errorf("failed to read property: %v", err)
		}

		properties, ok := props[property.PropID()]
		if ok {
			props[property.PropID()] = append(properties, property)
		} else {
			props[property.PropID()] = []Property{property}
		}
	}
	return props, nil
}

func readByteProp(reader io.Reader) (PropertyPayload, error) {
	var buf [1]byte
	if _, err := io.ReadFull(reader, buf[:1]); err != nil {
		return nil, fmt.Errorf("failed to read byte property: %v", err)
	}
	return BytePropPayload(buf[0]), nil
}

func readInt32Prop(reader io.Reader) (PropertyPayload, error) {
	val, err := types.ReadUInt32(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read four byte integer property: %v", err)
	}
	return Int32PropPayload(val), nil
}

func readInt16Prop(reader io.Reader) (PropertyPayload, error) {
	val, err := types.ReadUInt16(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read two byte integer property: %v", err)
	}
	return Int16PropPayload(val), nil
}

func readStringProp(reader io.Reader) (PropertyPayload, error) {
	val, err := types.ReadString(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read string property: %v", err)
	}
	return StringPropPayload(val), nil
}

func readKeyValueProp(reader io.Reader) (PropertyPayload, error) {
	key, err := types.ReadString(reader)
	if err != nil {
		return KeyValuePropPayload{}, fmt.Errorf("failed to read string pair property: failed to read key: %v", err)
	}
	value, err := types.ReadString(reader)
	if err != nil {
		return KeyValuePropPayload{}, fmt.Errorf("failed to read string pair property: failed to read value: %v", err)
	}
	return KeyValuePropPayload{key: key, value: value}, nil
}

func readVarIntProp(reader io.Reader) (PropertyPayload, error) {
	val, err := types.ReadVarInt(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read variable length int property: %v", err)
	}
	return VarIntPropPayload(val), nil
}

func readBinaryProp(reader io.Reader) (PropertyPayload, error) {
	val, err := types.ReadBinary(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read string pair property: %v", err)
	}
	return BinaryPropPayload(val), nil
}
