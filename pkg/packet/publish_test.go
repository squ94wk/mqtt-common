package packet

import (
	"bytes"
	"io"
	"testing"

	"github.com/go-test/deep"
	"github.com/squ94wk/mqtt-common/internal/help"
)

func TestReadPublish(t *testing.T) {
	type args struct {
		reader io.Reader
	}

	tests := []struct {
		name    string
		args    args
		want    *Publish
		wantErr bool
	}{
		{
			name: "publish1",
			args: args{
				reader: bytes.NewReader(publish1Bin.Bytes())},
			want: &publish1,
		},

		{
			name: "publish2",
			args: args{
				reader: bytes.NewReader(publish2Bin.Bytes())},
			want: &publish2,
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

func TestWritePublish(t *testing.T) {
	tests := []struct {
		name       string
		pkt        Publish
		wantWriter help.ByteSequence
		wantErr    bool
	}{
		{name: "publish1", pkt: publish1, wantWriter: publish1Bin},
		{name: "publish2", pkt: publish2, wantWriter: publish2Bin},
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
