package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

//Subscribe defines the subscribe control packet.
type Subscribe struct {
	PacketID uint16
	Props    Properties
	Filters  []SubscriptionFilter
}

//SubscriptionFilter defines all features of a single subscription.
type SubscriptionFilter struct {
	Filter            string
	MaxQoS            byte
	NoLocal           bool
	RetainAsPublished bool
	RetainHandling    byte
}

//RetainHandling* define readable aliases for three possible values for the retain handling.
const (
	RetainHandlingAlways       byte = 0
	RetainHandlingIfNotPresent byte = 1
	RetainHandlingNever        byte = 2
)

//WriteTo writes the subscribe control packet to writer according to the mqtt protocol.
func (s Subscribe) WriteTo(writer io.Writer) (int64, error) {
	var n int64
	// 3.8.1 Fixed header
	firstByte := byte(SUBSCRIBE)<<4 | 2
	n1, err := writer.Write([]byte{firstByte})
	n += int64(n1)
	if err != nil {
		return n, fmt.Errorf("failed to write subscribe packet: failed to write fixed header: %v", err)
	}

	//3.8.2 Variable header
	var remainingLength = types.UInt16Size // packetID
	remainingLength += s.Props.size()
	for _, filter := range s.Filters {
		remainingLength += types.StringSize(filter.Filter) + 1
	}
	n2, err := types.WriteVarIntTo(writer, remainingLength)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write subscribe packet: failed to write packet length: %v", err)
	}

	n4, err := types.WriteUInt16To(writer, s.PacketID)
	n += n4
	if err != nil {
		return n, fmt.Errorf("failed to write subscribe packet: failed to write packetID: %v", err)
	}

	n5, err := s.Props.WriteTo(writer)
	n += n5
	if err != nil {
		return n, fmt.Errorf("failed to write subscribe packet: failed to write properties: %v", err)
	}

	for _, filter := range s.Filters {
		n6, err := types.WriteStringTo(writer, filter.Filter)
		n += n6
		if err != nil {
			return n, fmt.Errorf("failed to write subscribe packet: failed to write subscribe filter: %v", err)
		}

		var options byte
		options |= filter.MaxQoS
		if filter.NoLocal {
			options |= 1 << 2
		}
		if filter.RetainAsPublished {
			options |= 1 << 3
		}
		options |= filter.RetainHandling << 4
		n7, err := writer.Write([]byte{options})
		n += int64(n7)
		if err != nil {
			return n, fmt.Errorf("failed to write subscribe packet: failed to write properties: %v", err)
		}
	}

	return n, nil
}

func readSubscribe(reader io.Reader, subscribe *Subscribe) error {
	// 3.8.2 Variable header
	// 3.8.2.1 Subscribe packet ID
	packetID, err := types.ReadUInt16(reader)
	if err != nil {
		return fmt.Errorf("failed to read subscribe packet: failed to read packet ID: %v", err)
	}
	subscribe.PacketID = packetID

	// 3.8.2.2 Subscribe properties
	props, err := readProperties(reader)
	if err != nil {
		return fmt.Errorf("failed to read subscribe packet: failed to read properties: %v", err)
	}
	subscribe.Props = props

	// 3.8.3 Payload
	var filters []SubscriptionFilter
	for reader.(*io.LimitedReader).N > 0 {
		filter, err := types.ReadString(reader)
		if err != nil {
			return fmt.Errorf("failed to read subscribe packet: failed to read filter: %v", err)
		}

		var buf [1]byte
		_, err = io.ReadFull(reader, buf[:])
		if err != nil {
			return fmt.Errorf("failed to read subscribe packet: failed to read subscription options: %v", err)
		}
		options := buf[0]

		maxQoS := options & 3
		noLocal := options&(1<<2) > 0
		retainAsPublished := options&(1<<3) > 0
		retainHandling := options & (3 << 4) >> 4

		filters = append(filters, SubscriptionFilter{
			Filter:            filter,
			MaxQoS:            maxQoS,
			NoLocal:           noLocal,
			RetainAsPublished: retainAsPublished,
			RetainHandling:    retainHandling,
		})
	}
	subscribe.Filters = filters
	return nil
}
