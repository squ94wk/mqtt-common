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

		{
			name: "publish3",
			args: args{
				reader: bytes.NewReader(publish3Bin.Bytes())},
			want: &publish3,
		},

		{
			name: "publish4",
			args: args{
				reader: bytes.NewReader(publish4Bin.Bytes())},
			want: &publish4,
		},

		{
			name: "publish5",
			args: args{
				reader: bytes.NewReader(publish5Bin.Bytes())},
			want: &publish5,
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
		{name: "publish3", pkt: publish3, wantWriter: publish3Bin},
		{name: "publish4", pkt: publish4, wantWriter: publish4Bin},
		{name: "publish5", pkt: publish5, wantWriter: publish5Bin},
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
