package packet

import (
	"bytes"
	"io"
	"testing"

	"github.com/go-test/deep"
	"github.com/squ94wk/mqtt-common/internal/help"
)

func TestReadConnack(t *testing.T) {
	type args struct {
		reader io.Reader
		header Header
	}
	tests := []struct {
		name    string
		args    args
		want    Connack
		wantErr bool
	}{
		//{
		//name: "Invalid flags => err", args: args{
		//reader: bytes.NewReader(help.Concat(
		//[]byte{2, byte(Success)},
		//)),
		//},
		//wantErr: true,
		//},
		//{
		//name: "connack1",
		//args: args{
		//reader: bytes.NewReader(connack1Bin.Bytes()),
		//},
		//want: connack1,
		//},
		//{
		//name: "connack2",
		//args: args{
		//reader: bytes.NewReader(connack2Bin.Bytes()),
		//},
		//want: connack2,
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var header Header
			if err := ReadHeader(tt.args.reader, &header); err != nil {
				if !tt.wantErr {
					t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			var connack Connack
			if err := ReadConnack(tt.args.reader, &connack, header); err != nil {
				if !tt.wantErr {
					t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if diff := deep.Equal(connack, tt.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestConnackWrite(t *testing.T) {
	tests := []struct {
		name       string
		connack    Connack
		wantWriter help.ByteSequence
		wantErr    bool
	}{
		{
			name:       "connack1",
			connack:    connack1,
			wantWriter: connack1Bin,
		},
		{
			name:       "connack2",
			connack:    connack2,
			wantWriter: connack2Bin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			err := tt.connack.Write(writer)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("pkt.Write() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			got := writer.Bytes()
			if diff := help.Match(tt.wantWriter, got); diff != nil {
				t.Error(diff)
			}
		})
	}
}

var (
	connack1Bin = help.NewByteSequence(
		help.IN_ORDER,
		help.NewByteSegment(
			[]byte{byte(CONNACK) << 4, 17},
			//variable header
			[]byte{1, byte(Success)},
			//props length
			[]byte{14},
		),
		//props
		help.NewByteSequence(
			help.ANY_ORDER,
			help.NewByteSegment([]byte{byte(AssignedClientIdentifier), 0, 6, 'c', 'l', 'i', 'e', 'n', 't'}),
			help.NewByteSegment([]byte{byte(MaximumPacketSize), 0, 1, 0, 0}),
		),
	)

	connack2Bin = help.NewByteSequence(
		help.IN_ORDER,
		help.NewByteSegment(
			[]byte{byte(CONNACK) << 4, 14},
			//variable header
			[]byte{0, byte(Success)},
			//props length
			[]byte{11},
		),
		//props
		help.NewByteSequence(
			help.ANY_ORDER,
			help.NewByteSegment([]byte{byte(AssignedClientIdentifier), 0, 6, 'c', 'l', 'i', 'e', 'n', 't'}),
			help.NewByteSegment([]byte{byte(MaximumQoS), 1}),
		),
	)

	connack1 = Connack{
		connectReason:  Success,
		sessionPresent: true,
		props: BuildProps(
			NewStringProp(AssignedClientIdentifier, "client"),
			NewInt32Prop(MaximumPacketSize, 1<<16),
		),
	}

	connack2 = Connack{
		connectReason:  Success,
		sessionPresent: false,
		props: BuildProps(
			NewStringProp(AssignedClientIdentifier, "client"),
			NewByteProp(MaximumQoS, byte(Qos1)),
		),
	}
)
