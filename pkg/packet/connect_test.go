package packet

import (
	"bytes"
	"io"
	"testing"

	"github.com/go-test/deep"
	"github.com/squ94wk/mqtt-common/internal/help"
	"github.com/squ94wk/mqtt-common/pkg/topic"
)

func TestReadConnect(t *testing.T) {
	type args struct {
		reader io.Reader
	}

	tests := []struct {
		name    string
		args    args
		want    *Connect
		wantErr bool
	}{
		{
			name: "protocol version 3 => err",
			args: args{
				reader: bytes.NewReader(help.Concat(
					[]byte{1 << 4, 18},
					[]byte{0, 4, 'M', 'Q', 'T', 'T', 3, 2, 0, 10},
					[]byte{5, 17, 0, 0, 0, 17},
					[]byte{0, 0}))},
			wantErr: true,
		},

		{
			name: "connect1",
			args: args{
				reader: bytes.NewReader(connect1Bin.Bytes())},
			want: &connect1,
		},

		{
			name: "connect2",
			args: args{
				reader: bytes.NewReader(connect2Bin.Bytes())},
			want: &connect2,
		},

		{
			name: "connect3",
			args: args{
				reader: bytes.NewReader(connect3Bin.Bytes())},
			want: &connect3,
		},

		{
			name: "connect4",
			args: args{
				reader: bytes.NewReader(connect4Bin.Bytes())},
			want: &connect4,
		},

		{
			name: "connect5",
			args: args{
				reader: bytes.NewReader(connect5Bin.Bytes())},
			want: &connect5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkt, err := ReadPacket(tt.args.reader)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if diff := deep.Equal(tt.want, pkt); diff != nil {
				t.Error(diff)
				return
			}
		})
	}
}

func TestWriteConnect(t *testing.T) {
	tests := []struct {
		name       string
		pkt        Connect
		wantWriter help.ByteSequence
		wantErr    bool
	}{
		{name: "connect1", pkt: connect1, wantWriter: connect1Bin},
		{name: "connect2", pkt: connect2, wantWriter: connect2Bin},
		{name: "connect3", pkt: connect3, wantWriter: connect3Bin},
		{name: "connect4", pkt: connect4, wantWriter: connect4Bin},
		{name: "connect5", pkt: connect5, wantWriter: connect5Bin},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := tt.pkt.WriteTo(writer); (err != nil) != tt.wantErr {
				t.Errorf("pkt.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotWriter := writer.Bytes()
			//fmt.Println(tt.wantWriter)
			//fmt.Println(gotWriter)
			if diff := help.Match(tt.wantWriter, gotWriter); diff != nil {
				t.Error(diff)
			}
		})
	}
}

var (
	connect1Bin = help.NewByteSegment(
		[]byte{byte(CONNECT) << 4, 18},
		[]byte{0, 4, 'M', 'Q', 'T', 'T', 5, 2, 0, 10},
		[]byte{5, byte(SessionExpiryInterval), 0, 0, 0, 10},
		[]byte{0, 0},
	)

	connect2Bin = help.NewByteSequence(
		help.IN_ORDER,
		help.NewByteSegment(
			[]byte{byte(CONNECT) << 4, 46},
			[]byte{0, 4, 'M', 'Q', 'T', 'T', 5},
			[]byte{2},
			[]byte{0, 100},
			[]byte{33}),
		help.NewByteSequence(
			help.ANY_ORDER,
			help.NewByteSegment([]byte{byte(SessionExpiryInterval), 0, 0, 0, 100}),
			help.NewByteSegment([]byte{byte(UserProperty), 0, 3, 'k', 'e', 'y', 0, 5, 'v', 'a', 'l', 'u', 'e'}),
			help.NewByteSegment([]byte{byte(UserProperty), 0, 4, 'k', 'e', 'y', '2', 0, 6, 'v', 'a', 'l', 'u', 'e', '2'}),
		),
		help.NewByteSegment(
			[]byte{0, 0},
		),
	)

	connect3Bin = help.NewByteSequence(
		help.IN_ORDER,
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
			help.ANY_ORDER,
			// //session expiry interval
			help.NewByteSegment([]byte{byte(SessionExpiryInterval), 0, 0, 0, 100}),
			// //user properties
			help.NewByteSegment([]byte{byte(UserProperty), 0, 3, 'k', 'e', 'y', 0, 5, 'v', 'a', 'l', 'u', 'e'}),
			help.NewByteSegment([]byte{byte(UserProperty), 0, 4, 'k', 'e', 'y', '2', 0, 6, 'v', 'a', 'l', 'u', 'e', '2'}),
		),
		help.NewByteSegment(
			//clientId
			[]byte{0, 0},
		),
	)

	connect4Bin = help.NewByteSequence(
		help.IN_ORDER,
		help.NewByteSegment(
			[]byte{byte(CONNECT) << 4, 113},
			//fixed header
			[]byte{0, 4, 'M', 'Q', 'T', 'T', 5},
			//flags
			[]byte{244},
			//keep alive
			[]byte{0, 100},
			//properties
			// //len
			[]byte{33},
		),
		help.NewByteSequence(
			help.ANY_ORDER,
			// //session expiry interval
			help.NewByteSegment([]byte{byte(SessionExpiryInterval), 0, 0, 0, 100}),
			// //user properties
			help.NewByteSegment([]byte{byte(UserProperty), 0, 3, 'k', 'e', 'y', 0, 5, 'v', 'a', 'l', 'u', 'e'}),
			help.NewByteSegment([]byte{byte(UserProperty), 0, 4, 'k', 'e', 'y', '2', 0, 6, 'v', 'a', 'l', 'u', 'e', '2'}),
		),
		help.NewByteSegment(
			//clientId
			[]byte{0, 8, 'c', 'l', 'i', 'e', 'n', 't', 'I', 'd'},
			//will properties
			// //len
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
		help.IN_ORDER,
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
			//clientId
			[]byte{0, 0},
		),
	)

	connect1 = NewConnect(
		true,
		10,
		BuildProps(
			NewInt32Prop(SessionExpiryInterval, 10),
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
		BuildProps(
			NewInt32Prop(SessionExpiryInterval, 100),
			NewKeyValueProp(UserProperty, "key", "value"),
			NewKeyValueProp(UserProperty, "key2", "value2"),
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
		BuildProps(
			NewInt32Prop(SessionExpiryInterval, 100),
			NewKeyValueProp(UserProperty, "key", "value"),
			NewKeyValueProp(UserProperty, "key2", "value2"),
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
		BuildProps(
			NewInt32Prop(SessionExpiryInterval, 100),
			NewKeyValueProp(UserProperty, "key", "value"),
			NewKeyValueProp(UserProperty, "key2", "value2"),
		),
		"clientId",
		true,
		Qos2,
		BuildProps(
			NewKeyValueProp(UserProperty, "willKey", "willValue"),
		),
		topic.Topic("/will/topic"),
		[]byte("willPayload"),
		"user",
		[]byte("pwd"),
	)

	connect5 = NewConnect(
		true,
		10,
		BuildProps(
			NewInt32Prop(SessionExpiryInterval, 10),
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
)
