package packet

import (
	"bytes"
	"io"
	"testing"

	"github.com/go-test/deep"
	"github.com/squ94wk/mqtt-common/internal/help"
)

func TestReadDisconnect(t *testing.T) {
	type args struct {
		reader io.Reader
		header header
	}
	tests := []struct {
		name    string
		args    args
		want    *Disconnect
		wantErr bool
	}{
		{
			name: "disconnect1",
			args: args{
				reader: bytes.NewReader(disconnect1Bin.Bytes()),
			},
			want: &disconnect1,
		},

		{
			name: "disconnect2",
			args: args{
				reader: bytes.NewReader(disconnect2Bin.Bytes()),
			},
			want: &disconnect2,
		},

		{
			name: "disconnect3",
			args: args{
				reader: bytes.NewReader(disconnect3Bin.Bytes()),
			},
			want: &disconnect3,
		},

		{
			name: "disconnect4",
			args: args{
				reader: bytes.NewReader(disconnect4Bin.Bytes()),
			},
			want: &disconnect4,
		},

		{
			name: "disconnect5",
			args: args{
				reader: bytes.NewReader(disconnect5Bin.Bytes()),
			},
			want: &disconnect5,
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

func TestDisconnectWriteTo(t *testing.T) {
	tests := []struct {
		name       string
		disconnect Disconnect
		wantWriter help.ByteSequence
		wantErr    bool
	}{
		{
			name:       "disconnect1",
			disconnect: disconnect1,
			wantWriter: disconnect1Bin,
		},

		{
			name:       "disconnect2",
			disconnect: disconnect2,
			wantWriter: disconnect2Bin,
		},

		{
			name:       "disconnect3",
			disconnect: disconnect3,
			//optimization: see 3.14.2
			wantWriter: disconnect5Bin,
		},

		{
			name:       "disconnect4",
			disconnect: disconnect4,
			//optimization: see 3.14.2
			wantWriter: disconnect5Bin,
		},

		{
			name:       "disconnect5",
			disconnect: disconnect5,
			wantWriter: disconnect5Bin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			_, err := tt.disconnect.WriteTo(writer)
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
