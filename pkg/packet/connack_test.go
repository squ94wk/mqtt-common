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
		header header
	}
	tests := []struct {
		name    string
		args    args
		want    *Connack
		wantErr bool
	}{
		{
			name: "Invalid flags => err", args: args{
				reader: bytes.NewReader(help.Concat(
					[]byte{2, byte(Success)},
				)),
			},
			wantErr: true,
		},

		{
			name: "connack1",
			args: args{
				reader: bytes.NewReader(connack1Bin.Bytes()),
			},
			want: &connack1,
		},

		{
			name: "connack2",
			args: args{
				reader: bytes.NewReader(connack2Bin.Bytes()),
			},
			want: &connack2,
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
			}
		})
	}
}

func TestConnackWriteTo(t *testing.T) {
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
			err := tt.connack.WriteTo(writer)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("pkt.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
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
