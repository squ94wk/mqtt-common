package packet

import (
	"bytes"
	"io"
	"testing"

	"github.com/go-test/deep"
	"github.com/squ94wk/mqtt-common/internal/help"
)

func TestReadSuback(t *testing.T) {
	type args struct {
		reader io.Reader
		header header
	}
	tests := []struct {
		name    string
		args    args
		want    *Suback
		wantErr bool
	}{
		{
			name: "suback1",
			args: args{
				reader: bytes.NewReader(suback1Bin.Bytes()),
			},
			want: &suback1,
		},

		{
			name: "suback2",
			args: args{
				reader: bytes.NewReader(suback2Bin.Bytes()),
			},
			want: &suback2,
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

func TestSubackWriteTo(t *testing.T) {
	tests := []struct {
		name       string
		suback     Suback
		wantWriter help.ByteSequence
		wantErr    bool
	}{
		{
			name:       "suback1",
			suback:     suback1,
			wantWriter: suback1Bin,
		},

		{
			name:       "suback2",
			suback:     suback2,
			wantWriter: suback2Bin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			_, err := tt.suback.WriteTo(writer)
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
