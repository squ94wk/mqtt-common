package packet

import (
	"fmt"
	"io"

	"github.com/squ94wk/mqtt-common/internal/types"
)

//Subscribe defines the subscribe control packet.
type Subscribe struct {
	packetID uint16
	props    Properties
	filters  []SubscriptionFilter
}

//SubscriptionFilter defines all features of a single subscription.
type SubscriptionFilter struct {
	filter            string
	maxQoS            byte
	noLocal           bool
	retainAsPublished bool
	retainHandling    byte
}

//RetainHandling* define readable aliases for three possible values for the retain handling.
const (
	RetainHandlingAlways       byte = 0
	RetainHandlingIfNotPresent byte = 1
	RetainHandlingNever        byte = 2
)

//NewSubscriptionFilter is the constructor for the SubscriptionFilter type.
func NewSubscriptionFilter(filter string, maxQoS byte, noLocal bool, retainAsPublished bool, retainHandling byte) SubscriptionFilter {
	return SubscriptionFilter{
		filter:            filter,
		maxQoS:            maxQoS,
		noLocal:           noLocal,
		retainAsPublished: retainAsPublished,
		retainHandling:    retainHandling,
	}
}

//PacketID returns the value of the subscribe control packet.
func (s Subscribe) PacketID() uint16 {
	return s.packetID
}

//SetPacketID sets the value of the subscribe control packet.
func (s *Subscribe) SetPacketID(packetID uint16) {
	s.packetID = packetID
}

//Props returns the properties of the subscribe control packet.
func (s Subscribe) Props() Properties {
	return s.props
}

//SetProps replaces the properties of the subscribe control packet.
func (s *Subscribe) SetProps(props map[uint32][]Property) {
	s.props = props
}

//SubscriptionFilters returns the filters of the subscribe control packet.
func (s Subscribe) SubscriptionFilters() []SubscriptionFilter {
	return s.filters
}

//SetSubscriptionFilters sets the filters of the subscribe control packet.
func (s *Subscribe) SetSubscriptionFilters(filters []SubscriptionFilter) {
	s.filters = filters
}

//Filter returns the topic filter of the subscription.
func (f SubscriptionFilter) Filter() string {
	return f.filter
}

//SetFilter sets the topic filter of the subscription.
func (f *SubscriptionFilter) SetFilter(filter string) {
	f.filter = filter
}

//MaxQoS returns the maximum supported quality of service level of the subscription.
func (f SubscriptionFilter) MaxQoS() byte {
	return f.maxQoS
}

//SetMaxQoS sets the topic filter of the subscription.
func (f *SubscriptionFilter) SetMaxQoS(qos byte) {
	f.maxQoS = qos
}

//NoLocal returns if own messages that match this subscription are published back to the client.
func (f SubscriptionFilter) NoLocal() bool {
	return f.noLocal
}

//SetNoLocal sets the "no local" option of the subscription.
func (f *SubscriptionFilter) SetNoLocal(noLocal bool) {
	f.noLocal = noLocal
}

//RetainAsPublished returns if the retain flag will not be changed through the server.
func (f SubscriptionFilter) RetainAsPublished() bool {
	return f.retainAsPublished
}

//SetRetainAsPublished sets the "retain as published" option of the subscription.
func (f *SubscriptionFilter) SetRetainAsPublished(retainAsPublished bool) {
	f.retainAsPublished = retainAsPublished
}

//RetainHandling returns if messages matching this subscriptions should be retained (and replayed if the session is continued).
func (f SubscriptionFilter) RetainHandling() byte {
	return f.retainHandling
}

//SetRetainHandling sets the retain handling option of the subscription.
func (f *SubscriptionFilter) SetRetainHandling(retainHandling byte) {
	f.retainHandling = retainHandling
}

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
	remainingLength += s.props.size()
	for _, filter := range s.filters {
		remainingLength += types.StringSize(filter.filter) + 1
	}
	n2, err := types.WriteVarIntTo(writer, remainingLength)
	n += n2
	if err != nil {
		return n, fmt.Errorf("failed to write subscribe packet: failed to write packet length: %v", err)
	}

	n4, err := types.WriteUInt16To(writer, s.packetID)
	n += n4
	if err != nil {
		return n, fmt.Errorf("failed to write subscribe packet: failed to write packetID: %v", err)
	}

	n5, err := s.props.WriteTo(writer)
	n += n5
	if err != nil {
		return n, fmt.Errorf("failed to write subscribe packet: failed to write properties: %v", err)
	}

	for _, filter := range s.filters {
		n6, err := types.WriteStringTo(writer, filter.filter)
		n += n6
		if err != nil {
			return n, fmt.Errorf("failed to write subscribe packet: failed to write subscribe filter: %v", err)
		}

		var options byte
		options |= filter.maxQoS
		if filter.noLocal {
			options |= 1 << 2
		}
		if filter.retainAsPublished {
			options |= 1 << 3
		}
		options |= filter.retainHandling << 4
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
	subscribe.SetPacketID(packetID)

	// 3.8.2.2 Subscribe properties
	props, err := readProperties(reader)
	if err != nil {
		return fmt.Errorf("failed to read subscribe packet: failed to read properties: %v", err)
	}
	subscribe.props = props

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
			filter:            filter,
			maxQoS:            maxQoS,
			noLocal:           noLocal,
			retainAsPublished: retainAsPublished,
			retainHandling:    retainHandling,
		})
	}
	subscribe.SetSubscriptionFilters(filters)
	return nil
}
