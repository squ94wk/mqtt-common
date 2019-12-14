package packet

import (
	"github.com/squ94wk/mqtt-common/internal/help"
)

var (
	connect1Bin = help.NewByteSegment(
		[]byte{byte(CONNECT) << 4, 18},
		[]byte{0, 4, 'M', 'Q', 'T', 'T', 5, 2, 0, 10},
		[]byte{5, byte(SessionExpiryInterval), 0, 0, 0, 10},
		[]byte{0, 0},
	)

	connect2Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(CONNECT) << 4, 46},
			[]byte{0, 4, 'M', 'Q', 'T', 'T', 5},
			[]byte{2},
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
			[]byte{2},
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
			[]byte{244},
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
			[]byte{2},
			//keep alive
			[]byte{0, 10},
			//properties (len + session expiry interval)
			[]byte{5, byte(SessionExpiryInterval), 0, 0, 0, 10},
			//clientID
			[]byte{0, 0},
		),
	)

	connect1 = NewConnect(
		true,
		10,
		NewProperties(
			Property{propID: SessionExpiryInterval, payload:Int32PropPayload(10)},
		),
		"",
		false,
		Qos0,
		nil,
		"",
		nil,
		"",
		nil,
	)

	connect2 = NewConnect(
		true,
		100,
		NewProperties(
			Property{propID:SessionExpiryInterval, payload: Int32PropPayload(100)},
			Property{propID:UserProperty, payload: NewKeyValuePropPayload("key", "value")},
			Property{propID:UserProperty, payload: NewKeyValuePropPayload("key2", "value2")},
		),
		"",
		false,
		Qos0,
		nil,
		"",
		nil,
		"",
		nil,
	)

	connect3 = NewConnect(
		true,
		100,
		NewProperties(
			Property{propID:SessionExpiryInterval, payload: Int32PropPayload(100)},
			Property{propID:UserProperty, payload: NewKeyValuePropPayload("key", "value")},
			Property{propID:UserProperty, payload: NewKeyValuePropPayload("key2", "value2")},
		),
		"",
		false,
		Qos0,
		nil,
		"",
		nil,
		"",
		nil,
	)

	connect4 = NewConnect(
		false,
		100,
		NewProperties(
			Property{propID:SessionExpiryInterval, payload: Int32PropPayload(100)},
			Property{propID:UserProperty, payload: NewKeyValuePropPayload("key", "value")},
			Property{propID:UserProperty, payload: NewKeyValuePropPayload("key2", "value2")},
		),
		"clientID",
		true,
		Qos2,
		NewProperties(
			Property{propID:UserProperty, payload: NewKeyValuePropPayload("willKey", "willValue")},
		),
		"/will/topic",
		[]byte("willPayload"),
		"user",
		[]byte("pwd"),
	)

	connect5 = NewConnect(
		true,
		10,
		NewProperties(
			Property{propID:SessionExpiryInterval, payload: Int32PropPayload(10)},
		),
		"",
		false,
		Qos0,
		nil,
		"",
		nil,
		"",
		nil,
	)

	connack1Bin = help.NewByteSequence(
		help.InOrder,
		help.NewByteSegment(
			[]byte{byte(CONNACK) << 4, 17},
			//variable header
			[]byte{1, byte(Success)},
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
			[]byte{0, byte(Success)},
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
		connectReason:  Success,
		sessionPresent: true,
		props: NewProperties(
			Property{propID:AssignedClientIdentifier, payload: StringPropPayload("client")},
			Property{propID:MaximumPacketSize, payload: Int32PropPayload(1<<16)},
		),
	}

	connack2 = Connack{
		connectReason:  Success,
		sessionPresent: false,
		props: NewProperties(
			Property{propID:AssignedClientIdentifier, payload: StringPropPayload("client")},
			Property{propID:MaximumQoS, payload: BytePropPayload(Qos1)},
		),
	}
)
