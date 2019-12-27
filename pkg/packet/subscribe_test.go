package packet

import (
	"bytes"
	"io"
	"testing"

	"github.com/go-test/deep"
	"github.com/squ94wk/mqtt-common/internal/help"
)

func TestReadSubscribe(t *testing.T) {
	type args struct {
		reader io.Reader
		header header
	}
	tests := []struct {
		name    string
		args    args
		want    *Subscribe
		wantErr bool
	}{
		{
			name: "subscribe1",
			args: args{
				reader: bytes.NewReader(subscribe1Bin.Bytes()),
			},
			want: &subscribe1,
		},

		{
			name: "subscribe2",
			args: args{
				reader: bytes.NewReader(subscribe2Bin.Bytes()),
			},
			want: &subscribe2,
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

func TestSubscribeWriteTo(t *testing.T) {
	tests := []struct {
		name       string
		subscribe  Subscribe
		wantWriter help.ByteSequence
		wantErr    bool
	}{
		{
			name:       "subscribe1",
			subscribe:  subscribe1,
			wantWriter: subscribe1Bin,
		},

		{
			name:       "subscribe2",
			subscribe:  subscribe2,
			wantWriter: subscribe2Bin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			_, err := tt.subscribe.WriteTo(writer)
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
