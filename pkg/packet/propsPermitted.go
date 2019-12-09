package packet

var (
	allowedProps map[uint8]map[PropID]int
)

type entry struct {
	part  uint8
	times int
}

func init() {
	const (
		Connect uint8 = iota
		Connack
		Publish
		Puback
		Pubrec
		Pubrel
		Pubcomp
		Subscribe
		Suback
		Unsubscribe
		Unsuback
		Disconnect
		Auth
		Will
	)

	allowedProps = make(map[uint8]map[PropID]int)

	// 2.2.2.2
	compatible(PayloadFormatIndicator, entry{Publish, 1}, entry{Will, 1})
	compatible(MessageExpiryInterval, entry{Publish, 1}, entry{Will, 1})
	compatible(ContentType, entry{Publish, 1}, entry{Will, 1})
	compatible(ResponseTopic, entry{Publish, 1}, entry{Will, 1})
	compatible(CorrelationData, entry{Publish, 1}, entry{Will, 1})
	compatible(SubscriptionIdentifier, entry{Publish, 1}, entry{Subscribe, 1})
	compatible(SessionExpiryInterval, entry{Connect, 1}, entry{Connack, 1}, entry{Disconnect, 1})
	compatible(AssignedClientIdentifier, entry{Connack, 1})
	compatible(ServerKeepAlive, entry{Connack, 1})
	compatible(AuthenticationMethod, entry{Connect, 1}, entry{Connack, 1}, entry{Auth, 1})
	compatible(AuthenticationData, entry{Connect, 1}, entry{Connack, 1}, entry{Auth, 1})
	compatible(RequestProblemInformation, entry{Connect, 1})
	compatible(RequestResponseInformation, entry{Connect, 1})
	compatible(ResponseInformation, entry{Connack, 1})
	compatible(ServerReference, entry{Connack, 1}, entry{Disconnect, 1})
	compatible(ReasonString, entry{Connack, 1}, entry{Puback, 1}, entry{Pubrec, 1}, entry{Pubrel, 1}, entry{Pubcomp, 1}, entry{Suback, 1}, entry{Unsuback, 1}, entry{Disconnect, 1}, entry{Auth, 1})
	compatible(ReceiveMaximum, entry{Connect, 1}, entry{Connack, 1})
	compatible(TopicAliasMaximum, entry{Connect, 1}, entry{Connack, 1})
	compatible(TopicAlias, entry{Publish, 1})
	compatible(MaximumQoS, entry{Connack, 1})
	compatible(RetainAvailable, entry{Connack, 1})
	compatible(UserProperty, entry{Connect, 1}, entry{Connack, 1}, entry{Publish, 1}, entry{Will, 1}, entry{Puback, 1}, entry{Pubrec, 1}, entry{Pubrel, 1}, entry{Pubcomp, 1}, entry{Subscribe, 1}, entry{Suback, 1}, entry{Unsubscribe, 1}, entry{Unsuback, 1}, entry{Disconnect, 1}, entry{Auth, 1})
	compatible(MaximumPacketSize, entry{Connect, 1}, entry{Connack, 1})
	compatible(WildcardSubscriptionAvailable, entry{Connack, 1})
	compatible(SubscriptionIdentifierAvailable, entry{Connack, 1})
	compatible(SharedSubscriptionAvailable, entry{Connack, 1})
}

func compatible(id PropID, entries ...entry) {
	for _, entry := range entries {
		_, ok := allowedProps[entry.part]
		if !ok {
			allowedProps[entry.part] = make(map[PropID]int)
		}

		allowedProps[entry.part][id] = entry.times
	}
}

//Allowed returns the number of times a property with identifier id may be used in a part.
func Allowed(part uint8, id PropID) int {
	times, ok := allowedProps[part]
	if !ok {
		return 0
	}

	return times[id]
}
