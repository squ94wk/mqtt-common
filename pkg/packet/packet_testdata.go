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

	connect1 = Connect{
		keepAlive:  10,
		cleanStart: true,
		props: NewProperties(
			Property{propID: SessionExpiryInterval, payload: Int32PropPayload(10)},
		),
		payload: ConnectPayload{
			clientID:    "",
			willProps:   nil,
			willRetain:  false,
			willQoS:     Qos0,
			willTopic:   "",
			willPayload: nil,
			username:    "",
			password:    nil,
		},
	}

	connect2 = Connect{
		keepAlive:  100,
		cleanStart: true,
		props: NewProperties(
			Property{propID: SessionExpiryInterval, payload: Int32PropPayload(100)},
			Property{propID: UserProperty, payload: NewKeyValuePropPayload("key", "value")},
			Property{propID: UserProperty, payload: NewKeyValuePropPayload("key2", "value2")},
		),
		payload: ConnectPayload{
			clientID:    "",
			willProps:   nil,
			willRetain:  false,
			willQoS:     Qos0,
			willTopic:   "",
			willPayload: nil,
			username:    "",
			password:    nil,
		},
	}

	connect3 = Connect{
		keepAlive:  100,
		cleanStart: true,
		props: NewProperties(
			Property{propID: SessionExpiryInterval, payload: Int32PropPayload(100)},
			Property{propID: UserProperty, payload: NewKeyValuePropPayload("key", "value")},
			Property{propID: UserProperty, payload: NewKeyValuePropPayload("key2", "value2")},
		),
		payload: ConnectPayload{
			clientID:    "",
			willProps:   nil,
			willRetain:  false,
			willQoS:     Qos0,
			willTopic:   "",
			willPayload: nil,
			username:    "",
			password:    nil,
		},
	}

	connect4 = Connect{
		keepAlive:  100,
		cleanStart: false,
		props: NewProperties(
			Property{propID: SessionExpiryInterval, payload: Int32PropPayload(100)},
			Property{propID: UserProperty, payload: NewKeyValuePropPayload("key", "value")},
			Property{propID: UserProperty, payload: NewKeyValuePropPayload("key2", "value2")},
		),
		payload: ConnectPayload{
			clientID: "clientID",
			willProps: NewProperties(
				Property{propID: UserProperty, payload: NewKeyValuePropPayload("willKey", "willValue")},
			),
			willRetain:  true,
			willQoS:     Qos2,
			willTopic:   "/will/topic",
			willPayload: []byte("willPayload"),
			username:    "user",
			password:    []byte("pwd"),
		},
	}

	connect5 = Connect{
		keepAlive:  10,
		cleanStart: true,
		props: NewProperties(
			Property{propID: SessionExpiryInterval, payload: Int32PropPayload(10)},
		),
		payload: ConnectPayload{
			clientID:    "",
			willProps:   nil,
			willRetain:  false,
			willQoS:     Qos0,
			willTopic:   "",
			willPayload: nil,
			username:    "",
			password:    nil,
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
		connectReason:  ConnectSuccess,
		sessionPresent: true,
		props: NewProperties(
			Property{propID: AssignedClientIdentifier, payload: StringPropPayload("client")},
			Property{propID: MaximumPacketSize, payload: Int32PropPayload(1 << 16)},
		),
	}

	connack2 = Connack{
		connectReason:  ConnectSuccess,
		sessionPresent: false,
		props: NewProperties(
			Property{propID: AssignedClientIdentifier, payload: StringPropPayload("client")},
			Property{propID: MaximumQoS, payload: BytePropPayload(Qos1)},
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
		reason: DisconnectImplementationSpecificError,
		props: NewProperties(
			Property{propID: ReasonString, payload: StringPropPayload("error")},
			Property{propID: SessionExpiryInterval, payload: Int32PropPayload(100)},
		),
	}

	disconnect2 = Disconnect{
		reason: DisconnectNormalDisconnection,
		props: NewProperties(
			Property{propID: UserProperty, payload: NewKeyValuePropPayload("key", "value")},
		),
	}

	disconnect3 = Disconnect{
		reason: DisconnectNormalDisconnection,
		props:  NewProperties(),
	}

	disconnect4 = Disconnect{
		reason: DisconnectNormalDisconnection,
		props:  NewProperties(),
	}

	disconnect5 = Disconnect{
		reason: DisconnectNormalDisconnection,
		props:  NewProperties(),
	}
)
