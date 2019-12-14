package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

//WriteTo writes a property to writer according to the mqtt protocol.
func (p Property) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	n1, err := types.WriteVarIntTo(writer, p.propID)
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write property: failed to write identifier: %v", err)
	}

	n2, err := p.payload.WriteTo(writer)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write property: failed to write payload: %v", err)
	}

	return n, nil
}

//WriteTo is an auxiliary function to write all properties from a map to writer.
func (p Properties) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	var propsSize uint32
	for propID, propsForID := range p {
		for _, prop := range propsForID {
			propsSize += types.VarIntSize(propID)
			propsSize += prop.payload.size()
		}
	}
	n1, err := types.WriteVarIntTo(writer, propsSize)
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write properties: failed to write size: %v", err)
	}

	for _, propsForID := range p {
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

//WriteTo writes the byte property to writer according to the mqtt protocol.
func (p BytePropPayload) WriteTo(writer io.Writer) (int64, error) {
	n, err := writer.Write([]byte{byte(p)})
	if err != nil {
		return int64(n), fmt.Errorf("failed to write byte property. failed to write payload: %v", err)
	}

	return int64(n), nil
}

//WriteTo writes the 32 bit integer property to writer according to the mqtt protocol.
func (p Int32PropPayload) WriteTo(writer io.Writer) (int64, error) {
	n, err := types.WriteUInt32To(writer, uint32(p))
	if err != nil {
		return n, fmt.Errorf("failed to write four byte property. failed to write payload: %v", err)
	}

	return n, nil
}

//WriteTo writes the 16 bit integer property to writer according to the mqtt protocol.
func (p Int16PropPayload) WriteTo(writer io.Writer) (int64, error) {
	n, err := types.WriteUInt16To(writer, uint16(p))
	if err != nil {
		return n, fmt.Errorf("failed to write two byte property. failed to write payload: %v", err)
	}

	return n, nil
}

//WriteTo writes the string property to writer according to the mqtt protocol.
func (p StringPropPayload) WriteTo(writer io.Writer) (int64, error) {
	n, err := types.WriteStringTo(writer, string(p))
	if err != nil {
		return n, fmt.Errorf("failed to write utf8 string property. failed to write payload: %v", err)
	}

	return n, nil
}

//WriteTo writes the key value property to writer according to the mqtt protocol.
func (p KeyValuePropPayload) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	n1, err := types.WriteStringTo(writer, p.key)
	n += n1
	if err != nil {
		return n, fmt.Errorf("failed to write string pair property. failed to write key: %v", err)
	}

	n2, err := types.WriteStringTo(writer, p.value)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write string pair property. failed to write value: %v", err)
	}

	return n, nil
}

//WriteTo writes the variable length integer property to writer according to the mqtt protocol.
func (p VarIntPropPayload) WriteTo(writer io.Writer) (int64, error) {
	n, err := types.WriteVarIntTo(writer, uint32(p))
	if err != nil {
		return n, fmt.Errorf("failed to write variable length integer property. failed to write payload: %v", err)
	}

	return n, nil
}

//WriteTo writes the binary property to writer according to the mqtt protocol.
func (p BinaryPropPayload) WriteTo(writer io.Writer) (int64, error) {
	n, err := types.WriteBinaryTo(writer, p)
	if err != nil {
		return n, fmt.Errorf("failed to write binary property. failed to write payload: %v", err)
	}

	return n, nil
}
