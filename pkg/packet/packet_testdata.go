package packet

import (
	"github.com/go-test/deep"
	"github.com/squ94wk/mqtt-common/internal/help"
	"github.com/squ94wk/mqtt-common/pkg/topic"
)

func init() {
	deep.CompareUnexportedFields = true
}

var (
	connect1Bin = help.NewByteSegment(
		[]byte{byte(CONNECT) << 4, 18},
		[]byte{0, 4, 'M', 'Q', 'T', 'T', 5, 1 << 1, 0, 10},
		[]byte{5, byte(SessionExpiryInterval), 0, 0, 0, 10},
		[]byte{0, 0},
	)

	connect2Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(CONNECT) << 4, 46},
			[]byte{0, 4, 'M', 'Q', 'T', 'T', 5},
			[]byte{1 << 1},
			[]byte{0, 100},
			[]byte{33}),
		help.NewByteSequence(
			help.AnyOrder,
			help.NewByteSegment([]byte{byte(SessionExpiryInterval), 0, 0, 0, 100}),
			help.NewByteSegment([]byte{byte(UserProperty), 0, 3, 'k', 'e', 'y', 0, 5, 'v', 'a', 'l', 'u', 'e'}),
			help.NewByteSegment([]byte{byte(UserProperty), 0, 4, 'k', 'e', 'y', '2', 0, 6, 'v', 'a', 'l', 'u', 'e', '2'}),
		),
		help.NewByteSegment(
			[]byte{0, 0},
		),
	)

	connect3Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			//fixed header
			[]byte{byte(CONNECT) << 4, 46},
			[]byte{0, 4, 'M', 'Q', 'T', 'T', 5},
			//flags
			[]byte{1 << 1},
			//keep alive
			[]byte{0, 100},
			//properties
			// // len
			[]byte{33},
		),
		help.NewByteSequence(
			help.AnyOrder,
			// //session expiry interval
			help.NewByteSegment([]byte{byte(SessionExpiryInterval), 0, 0, 0, 100}),
			// //user properties
			help.NewByteSegment([]byte{byte(UserProperty), 0, 3, 'k', 'e', 'y', 0, 5, 'v', 'a', 'l', 'u', 'e'}),
			help.NewByteSegment([]byte{byte(UserProperty), 0, 4, 'k', 'e', 'y', '2', 0, 6, 'v', 'a', 'l', 'u', 'e', '2'}),
		),
		help.NewByteSegment(
			//clientID
			[]byte{0, 0},
		),
	)

	connect4Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(CONNECT) << 4, 113},
			//fixed header
			[]byte{0, 4, 'M', 'Q', 'T', 'T', 5},
			//flags
			[]byte{1<<2 | 2<<3 | 1<<5 | 1<<6 | 1<<7},
			//keep alive
			[]byte{0, 100},
			//properties length
			[]byte{33},
		),
		help.NewByteSequence(
			help.AnyOrder,
			// //session expiry interval
			help.NewByteSegment([]byte{byte(SessionExpiryInterval), 0, 0, 0, 100}),
			// //user properties
			help.NewByteSegment([]byte{byte(UserProperty), 0, 3, 'k', 'e', 'y', 0, 5, 'v', 'a', 'l', 'u', 'e'}),
			help.NewByteSegment([]byte{byte(UserProperty), 0, 4, 'k', 'e', 'y', '2', 0, 6, 'v', 'a', 'l', 'u', 'e', '2'}),
		),
		help.NewByteSegment(
			//clientID
			[]byte{0, 8, 'c', 'l', 'i', 'e', 'n', 't', 'I', 'D'},
			//will properties length
			[]byte{21},
			// //user prop
			[]byte{byte(UserProperty), 0, 7, 'w', 'i', 'l', 'l', 'K', 'e', 'y', 0, 9, 'w', 'i', 'l', 'l', 'V', 'a', 'l', 'u', 'e'},
			// //will topic
			[]byte{0, 11, '/', 'w', 'i', 'l', 'l', '/', 't', 'o', 'p', 'i', 'c'},
			// //will payload
			[]byte{0, 11, 'w', 'i', 'l', 'l', 'P', 'a', 'y', 'l', 'o', 'a', 'd'},
			// //username
			[]byte{0, 4, 'u', 's', 'e', 'r'},
			// //password
			[]byte{0, 3, 'p', 'w', 'd'},
		),
	)

	connect5Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(CONNECT) << 4, 18},
			//fixed header
			[]byte{0, 4, 'M', 'Q', 'T', 'T', 5},
			//flags
			[]byte{1 << 1},
			//keep alive
			[]byte{0, 10},
			//properties (len + session expiry interval)
			[]byte{5, byte(SessionExpiryInterval), 0, 0, 0, 10},
			//clientID
			[]byte{0, 0},
		),
	)

	connect1 = Connect{
		KeepAlive:  10,
		CleanStart: true,
		Props: NewProperties(
			Property{PropID: SessionExpiryInterval, Payload: Int32PropPayload(10)},
		),
		Payload: ConnectPayload{
			ClientID:    "",
			WillProps:   nil,
			WillRetain:  false,
			WillQoS:     Qos0,
			WillTopic:   "",
			WillPayload: nil,
			Username:    "",
			Password:    nil,
		},
	}

	connect2 = Connect{
		KeepAlive:  100,
		CleanStart: true,
		Props: NewProperties(
			Property{PropID: SessionExpiryInterval, Payload: Int32PropPayload(100)},
			Property{PropID: UserProperty, Payload: KeyValuePropPayload{"key", "value"}},
			Property{PropID: UserProperty, Payload: KeyValuePropPayload{"key2", "value2"}},
		),
		Payload: ConnectPayload{
			ClientID:    "",
			WillProps:   nil,
			WillRetain:  false,
			WillQoS:     Qos0,
			WillTopic:   "",
			WillPayload: nil,
			Username:    "",
			Password:    nil,
		},
	}

	connect3 = Connect{
		KeepAlive:  100,
		CleanStart: true,
		Props: NewProperties(
			Property{PropID: SessionExpiryInterval, Payload: Int32PropPayload(100)},
			Property{PropID: UserProperty, Payload: KeyValuePropPayload{"key", "value"}},
			Property{PropID: UserProperty, Payload: KeyValuePropPayload{"key2", "value2"}},
		),
		Payload: ConnectPayload{
			ClientID:    "",
			WillProps:   nil,
			WillRetain:  false,
			WillQoS:     Qos0,
			WillTopic:   "",
			WillPayload: nil,
			Username:    "",
			Password:    nil,
		},
	}

	connect4 = Connect{
		KeepAlive:  100,
		CleanStart: false,
		Props: NewProperties(
			Property{PropID: SessionExpiryInterval, Payload: Int32PropPayload(100)},
			Property{PropID: UserProperty, Payload: KeyValuePropPayload{"key", "value"}},
			Property{PropID: UserProperty, Payload: KeyValuePropPayload{"key2", "value2"}},
		),
		Payload: ConnectPayload{
			ClientID: "clientID",
			WillProps: NewProperties(
				Property{PropID: UserProperty, Payload: KeyValuePropPayload{"willKey", "willValue"}},
			),
			WillRetain:  true,
			WillQoS:     Qos2,
			WillTopic:   "/will/topic",
			WillPayload: []byte("willPayload"),
			Username:    "user",
			Password:    []byte("pwd"),
		},
	}

	connect5 = Connect{
		KeepAlive:  10,
		CleanStart: true,
		Props: NewProperties(
			Property{PropID: SessionExpiryInterval, Payload: Int32PropPayload(10)},
		),
		Payload: ConnectPayload{
			ClientID:    "",
			WillProps:   nil,
			WillRetain:  false,
			WillQoS:     Qos0,
			WillTopic:   "",
			WillPayload: nil,
			Username:    "",
			Password:    nil,
		},
	}

	connack1Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(CONNACK) << 4, 17},
			//variable header
			[]byte{1, byte(ConnectSuccess)},
			//props length
			[]byte{14},
		),
		//props
		help.NewByteSequence(
			help.AnyOrder,
			help.NewByteSegment([]byte{byte(AssignedClientIdentifier), 0, 6, 'c', 'l', 'i', 'e', 'n', 't'}),
			help.NewByteSegment([]byte{byte(MaximumPacketSize), 0, 1, 0, 0}),
		),
	)

	connack2Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(CONNACK) << 4, 14},
			//variable header
			[]byte{0, byte(ConnectSuccess)},
			//props length
			[]byte{11},
		),
		//props
		help.NewByteSequence(
			help.AnyOrder,
			help.NewByteSegment([]byte{byte(AssignedClientIdentifier), 0, 6, 'c', 'l', 'i', 'e', 'n', 't'}),
			help.NewByteSegment([]byte{byte(MaximumQoS), 1}),
		),
	)

	connack1 = Connack{
		ConnectReason:  ConnectSuccess,
		SessionPresent: true,
		Props: NewProperties(
			Property{PropID: AssignedClientIdentifier, Payload: StringPropPayload("client")},
			Property{PropID: MaximumPacketSize, Payload: Int32PropPayload(1 << 16)},
		),
	}

	connack2 = Connack{
		ConnectReason:  ConnectSuccess,
		SessionPresent: false,
		Props: NewProperties(
			Property{PropID: AssignedClientIdentifier, Payload: StringPropPayload("client")},
			Property{PropID: MaximumQoS, Payload: BytePropPayload(Qos1)},
		),
	}

	disconnect1Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(DISCONNECT) << 4, 15},
			//variable header
			[]byte{byte(DisconnectImplementationSpecificError)},
			//props length
			[]byte{13},
		),
		//props
		help.NewByteSequence(
			help.AnyOrder,
			help.NewByteSegment([]byte{byte(ReasonString), 0, 5, 'e', 'r', 'r', 'o', 'r'}),
			help.NewByteSegment([]byte{byte(SessionExpiryInterval), 0, 0, 0, 100}),
		),
	)

	disconnect2Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(DISCONNECT) << 4, 15},
			//variable header
			[]byte{byte(DisconnectNormalDisconnection)},
			//props length
			[]byte{13},
		),
		//props
		help.NewByteSequence(
			help.AnyOrder,
			help.NewByteSegment([]byte{byte(UserProperty), 0, 3, 'k', 'e', 'y', 0, 5, 'v', 'a', 'l', 'u', 'e'}),
		),
	)

	disconnect3Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(DISCONNECT) << 4, 2},
			//variable header
			[]byte{byte(DisconnectNormalDisconnection)},
			//prop length
			[]byte{0},
		),
	)

	disconnect4Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(DISCONNECT) << 4, 1},
			//variable header
			[]byte{byte(DisconnectNormalDisconnection)},
		),
	)

	disconnect5Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(DISCONNECT) << 4, 0},
		),
	)

	disconnect1 = Disconnect{
		Reason: DisconnectImplementationSpecificError,
		Props: NewProperties(
			Property{PropID: ReasonString, Payload: StringPropPayload("error")},
			Property{PropID: SessionExpiryInterval, Payload: Int32PropPayload(100)},
		),
	}

	disconnect2 = Disconnect{
		Reason: DisconnectNormalDisconnection,
		Props: NewProperties(
			Property{PropID: UserProperty, Payload: KeyValuePropPayload{"key", "value"}},
		),
	}

	disconnect3 = Disconnect{
		Reason: DisconnectNormalDisconnection,
		Props:  NewProperties(),
	}

	disconnect4 = Disconnect{
		Reason: DisconnectNormalDisconnection,
		Props:  NewProperties(),
	}

	disconnect5 = Disconnect{
		Reason: DisconnectNormalDisconnection,
		Props:  NewProperties(),
	}

	subscribe1Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(SUBSCRIBE)<<4 | 2, 39},
			//variable header
			//packetID
			[]byte{0, 100},
			//props length
			[]byte{16},
		),
		//props
		help.NewByteSequence(
			help.AnyOrder,
			help.NewByteSegment([]byte{byte(SubscriptionIdentifier), 144, 78}),
			help.NewByteSegment([]byte{byte(UserProperty), 0, 3, 'k', 'e', 'y', 0, 5, 'v', 'a', 'l', 'u', 'e'}),
		),
		//subscriptions
		help.NewByteSequence(
			help.InOrder,
			help.NewByteSegment([]byte{0, 7, '/', 't', 'o', 'p', 'i', 'c', '1', 0}),
			help.NewByteSegment([]byte{0, 7, '/', 't', 'o', 'p', 'i', 'c', '2', 1 | 1<<2 | 1<<3 | 1<<4}),
		),
	)

	subscribe2Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(SUBSCRIBE)<<4 | 2, 16},
			//variable header
			//packetID
			[]byte{3, 232},
			//props length
			[]byte{3},
		),
		//props
		help.NewByteSequence(
			help.AnyOrder,
			help.NewByteSegment([]byte{byte(SubscriptionIdentifier), 232, 7}),
		),
		//subscriptions
		help.NewByteSequence(
			help.InOrder,
			help.NewByteSegment([]byte{0, 7, '/', 't', 'o', 'p', 'i', 'c', '3', 0 | 2 | 2<<4}),
		),
	)

	subscribe1 = Subscribe{
		PacketID: 100,
		Props: NewProperties(
			Property{PropID: UserProperty, Payload: KeyValuePropPayload{"key", "value"}},
			Property{PropID: SubscriptionIdentifier, Payload: VarIntPropPayload(10000)},
		),
		Filters: []SubscriptionFilter{
			{
				Filter:            "/topic1",
				MaxQoS:            Qos0,
				NoLocal:           false,
				RetainAsPublished: false,
				RetainHandling:    0,
			},
			{
				Filter:            "/topic2",
				MaxQoS:            Qos1,
				NoLocal:           true,
				RetainAsPublished: true,
				RetainHandling:    1,
			},
		},
	}

	subscribe2 = Subscribe{
		PacketID: 1000,
		Props: NewProperties(
			Property{PropID: SubscriptionIdentifier, Payload: VarIntPropPayload(1000)},
		),
		Filters: []SubscriptionFilter{
			{
				Filter:            "/topic3",
				MaxQoS:            Qos2,
				NoLocal:           false,
				RetainAsPublished: false,
				RetainHandling:    2,
			},
		},
	}

	suback1Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(SUBACK) << 4, 22},
			//variable header
			//packetID
			[]byte{0, 100},
			//props length
			[]byte{16},
		),
		//props
		help.NewByteSequence(
			help.AnyOrder,
			help.NewByteSegment([]byte{byte(SubscriptionIdentifier), 144, 78}),
			help.NewByteSegment([]byte{byte(UserProperty), 0, 3, 'k', 'e', 'y', 0, 5, 'v', 'a', 'l', 'u', 'e'}),
		),
		//reason codes
		help.NewByteSegment([]byte{
			byte(SubackGrantedQoS2),
			byte(SubackGrantedQoS0),
			byte(SubackQuotaExceeded),
		}),
	)

	suback2Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(SUBACK) << 4, 8},
			//variable header
			//packetID
			[]byte{3, 232},
			//props length
			[]byte{3},
		),
		//props
		help.NewByteSequence(
			help.AnyOrder,
			help.NewByteSegment([]byte{byte(SubscriptionIdentifier), 232, 7}),
		),
		//reason codes
		help.NewByteSegment([]byte{
			byte(SubackGrantedQoS0),
			byte(SubackImplementationSpecificError),
		}),
	)

	suback1 = Suback{
		PacketID: 100,
		Props: NewProperties(
			Property{PropID: UserProperty, Payload: KeyValuePropPayload{"key", "value"}},
			Property{PropID: SubscriptionIdentifier, Payload: VarIntPropPayload(10000)},
		),
		Reasons: []SubackReason{
			SubackGrantedQoS2,
			SubackGrantedQoS0,
			SubackQuotaExceeded,
		},
	}

	suback2 = Suback{
		PacketID: 1000,
		Props: NewProperties(
			Property{PropID: SubscriptionIdentifier, Payload: VarIntPropPayload(1000)},
		),
		Reasons: []SubackReason{
			SubackGrantedQoS0,
			SubackImplementationSpecificError,
		},
	}

	publish1Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(PUBLISH) << 4, 27},
			//variable header
			//topic name
			[]byte{0, 10},
			[]byte("device/abc"),
			//packetID
			[]byte{0, 100},
			//props length
			[]byte{5},
		),
		//props
		help.NewByteSegment([]byte{byte(MessageExpiryInterval), 0, 0, 0, 50}),
		//payload
		help.NewByteSegment([]byte("payload")),
	)

	publish2Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(PUBLISH)<<4 | 1<<3 | 2<<1 | 1, 34},
			//variable header
			//topic name
			[]byte{0, 15},
			[]byte("device/abc/temp"),
			//packetID
			[]byte{0, 100},
			//props length
			[]byte{7},
		),
		//props
		help.NewByteSequence(
			help.AnyOrder,
			help.NewByteSegment([]byte{byte(MessageExpiryInterval), 0, 0, 0, 50}),
			help.NewByteSegment([]byte{byte(PayloadFormatIndicator), 1}),
		),
		//payload
		help.NewByteSegment([]byte("payload")),
	)

	publish1 = Publish{
		Dup:      false,
		Qos:      Qos0,
		Retain:   false,
		Topic:    topic.Topic{Levels: []string{"device", "abc"}},
		PacketID: 100,
		Props:    NewProperties(
			Property{PropID: MessageExpiryInterval, Payload: Int32PropPayload(50)},
		),
		Payload:  []byte("payload"),
	}

	publish2 = Publish{
		Dup:      true,
		Qos:      Qos2,
		Retain:   true,
		Topic:    topic.Topic{Levels: []string{"device", "abc", "temp"}},
		PacketID: 100,
		Props:    NewProperties(
			Property{PropID: MessageExpiryInterval, Payload: Int32PropPayload(50)},
			Property{PropID: PayloadFormatIndicator, Payload: BytePropPayload(1)},
		),
		Payload:  []byte("payload"),
	}
)
