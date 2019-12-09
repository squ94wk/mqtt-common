package packet

import (
	"bytes"
	"io"
	"testing"

	"github.com/go-test/deep"
	"github.com/squ94wk/mqtt-common/internal/help"
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
			if _, err := tt.pkt.WriteTo(writer); (err != nil) != tt.wantErr {
				t.Errorf("pkt.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotWriter := writer.Bytes()
			if diff := help.Match(tt.wantWriter, gotWriter); diff != nil {
				t.Error(diff)
			}
		})
	}
}
