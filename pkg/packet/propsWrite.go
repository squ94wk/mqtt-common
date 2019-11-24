package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

func (prop ByteProp) Write(writer io.Writer) error {
	err := types.WriteVarInt(writer, uint32(prop.PropId()))
	if err != nil {
		return fmt.Errorf("failed to write byte property. failed to write identifier '%d': %v", prop.PropId(), err)
	}

	_, err = writer.Write([]byte{prop.value})
	if err != nil {
		return fmt.Errorf("failed to write byte property. failed to write payload: %v", err)
	}

	return nil
}

func (prop Int32Prop) Write(writer io.Writer) error {
	err := types.WriteVarInt(writer, uint32(prop.PropId()))
	if err != nil {
		return fmt.Errorf("failed to write four byte property. failed to write identifier '%d': %v", prop.PropId(), err)
	}

	err = types.WriteUInt32(writer, prop.value)
	if err != nil {
		return fmt.Errorf("failed to write four byte property. failed to write payload: %v", err)
	}

	return nil
}

func (prop Int16Prop) Write(writer io.Writer) error {
	err := types.WriteVarInt(writer, uint32(prop.PropId()))
	if err != nil {
		return fmt.Errorf("failed to write two byte property. failed to write identifier '%d': %v", prop.PropId(), err)
	}

	err = types.WriteUInt16(writer, prop.value)
	if err != nil {
		return fmt.Errorf("failed to write two byte property. failed to write payload: %v", err)
	}

	return nil
}

func (prop StringProp) Write(writer io.Writer) error {
	err := types.WriteVarInt(writer, uint32(prop.PropId()))
	if err != nil {
		return fmt.Errorf("failed to write utf8 string property. failed to write identifier '%d': %v", prop.PropId(), err)
	}

	err = types.WriteString(writer, prop.value)
	if err != nil {
		return fmt.Errorf("failed to write utf8 string property. failed to write payload: %v", err)
	}

	return nil
}

func (prop KeyValueProp) Write(writer io.Writer) error {
	err := types.WriteVarInt(writer, uint32(prop.PropId()))
	if err != nil {
		return fmt.Errorf("failed to write string pair property. failed to write identifier '%d': %v", prop.PropId(), err)
	}

	err = types.WriteStringPair(writer, prop.value)
	if err != nil {
		return fmt.Errorf("failed to write string pair property. failed to write payload: %v", err)
	}

	return nil
}

func (prop VarIntProp) Write(writer io.Writer) error {
	err := types.WriteVarInt(writer, uint32(prop.PropId()))
	if err != nil {
		return fmt.Errorf("failed to write variable length integer property. failed to write identifier '%d': %v", prop.PropId(), err)
	}

	err = types.WriteVarInt(writer, prop.value)
	if err != nil {
		return fmt.Errorf("failed to write variable length integer property. failed to write payload: %v", err)
	}

	return nil
}

func (prop BinaryProp) Write(writer io.Writer) error {
	err := types.WriteVarInt(writer, uint32(prop.PropId()))
	if err != nil {
		return fmt.Errorf("failed to write binary property. failed to write identifier '%d': %v", prop.PropId(), err)
	}

	err = types.WriteBinary(writer, prop.value)
	if err != nil {
		return fmt.Errorf("failed to write binary property. failed to write payload: %v", err)
	}

	return nil
}

func WriteProperties(writer io.Writer, props map[PropId][]Property) error {
	var propsSize uint32 = 0
	for _, propsForId := range props {
		for _, prop := range propsForId {
			propsSize += prop.size()
		}
	}
	if err := types.WriteVarInt(writer, uint32(propsSize)); err != nil {
		return fmt.Errorf("failed to write properties: failed to write size: %v", err)
	}

	for _, propsForId := range props {
		for _, prop := range propsForId {
			if err := prop.Write(writer); err != nil {
				return fmt.Errorf("failed to write properties: failed to write property with id '%d': %v", prop.PropId(), err)
			}
		}
	}
	return nil
}
