package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

//WriteTo writes the byte property to writer according to the mqtt protocol.
func (p ByteProp) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	n1, err := types.WriteVarIntTo(writer, uint32(p.PropID()))
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write byte property. failed to write identifier '%d': %v", p.PropID(), err)
	}

	n2, err := writer.Write([]byte{p.value})
	n += int64(n2)
	if err != nil {
		return n, fmt.Errorf("failed to write byte property. failed to write payload: %v", err)
	}

	return n, nil
}

//WriteTo writes the 32 bit integer property to writer according to the mqtt protocol.
func (p Int32Prop) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	n1, err := types.WriteVarIntTo(writer, uint32(p.PropID()))
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write four byte property. failed to write identifier '%d': %v", p.PropID(), err)
	}

	n2, err := types.WriteUInt32To(writer, p.value)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write four byte property. failed to write payload: %v", err)
	}

	return n, nil
}

//WriteTo writes the 16 bit integer property to writer according to the mqtt protocol.
func (p Int16Prop) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	n1, err := types.WriteVarIntTo(writer, uint32(p.PropID()))
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write two byte property. failed to write identifier '%d': %v", p.PropID(), err)
	}

	n2, err := types.WriteUInt16To(writer, p.value)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write two byte property. failed to write payload: %v", err)
	}

	return n, nil
}

//WriteTo writes the string property to writer according to the mqtt protocol.
func (p StringProp) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	n1, err := types.WriteVarIntTo(writer, uint32(p.PropID()))
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write utf8 string property. failed to write identifier '%d': %v", p.PropID(), err)
	}

	n2, err := types.WriteStringTo(writer, p.value)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write utf8 string property. failed to write payload: %v", err)
	}

	return n, nil
}

//WriteTo writes the key value property to writer according to the mqtt protocol.
func (p KeyValueProp) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	n1, err := types.WriteVarIntTo(writer, uint32(p.PropID()))
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write string pair property. failed to write identifier '%d': %v", p.PropID(), err)
	}

	n2, err := types.WriteStringTo(writer, p.key)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write string pair property. failed to write key: %v", err)
	}

	n3, err := types.WriteStringTo(writer, p.value)
	n += n3
	if err != nil {
		return n, fmt.Errorf("failed to write string pair property. failed to write value: %v", err)
	}

	return n, nil
}

//WriteTo writes the variable length integer property to writer according to the mqtt protocol.
func (p VarIntProp) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	n1, err := types.WriteVarIntTo(writer, uint32(p.PropID()))
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write variable length integer property. failed to write identifier '%d': %v", p.PropID(), err)
	}

	n2, err := types.WriteVarIntTo(writer, p.value)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write variable length integer property. failed to write payload: %v", err)
	}

	return n, nil
}

//WriteTo writes the binary property to writer according to the mqtt protocol.
func (p BinaryProp) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	n1, err := types.WriteVarIntTo(writer, uint32(p.PropID()))
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write binary property. failed to write identifier '%d': %v", p.PropID(), err)
	}

	n2, err := types.WriteBinaryTo(writer, p.value)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write binary property. failed to write payload: %v", err)
	}

	return n, nil
}

//WritePropsTo is an auxiliary function to write all properties from a map to writer.
func WritePropsTo(writer io.Writer, props map[PropID][]Property) (int64, error) {
	var n int64
	var propsSize uint32
	for _, propsForID := range props {
		for _, prop := range propsForID {
			propsSize += prop.size()
		}
	}
	n1, err := types.WriteVarIntTo(writer, propsSize)
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write properties: failed to write size: %v", err)
	}

	for _, propsForID := range props {
		for _, prop := range propsForID {
			n2, err := prop.WriteTo(writer)
			n += n2
			if err != nil {
				return n, fmt.Errorf("failed to write properties: failed to write property with id '%d': %v", prop.PropID(), err)
			}
		}
	}
	return n, nil
}
