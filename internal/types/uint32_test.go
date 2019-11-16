package types

import (
	"bytes"
	"io"
	"testing"
)

func TestWriteUInt32(t *testing.T) {
	type args struct {
		writer bytes.Buffer
		value  uint32
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{name: "0 => {0, 0, 0, 0}", args: args{bytes.Buffer{}, 0}, want: []byte{0, 0, 0, 0}},
		{name: "1 => {0, 0, 0, 1}", args: args{bytes.Buffer{}, 1}, want: []byte{0, 0, 0, 1}},
		{name: "256 => {0, 0, 1, 0}", args: args{bytes.Buffer{}, 256}, want: []byte{0, 0, 1, 0}},
		{name: "16777216 => {1, 0, 0, 0}", args: args{bytes.Buffer{}, 16777216}, want: []byte{1, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteUInt32(&tt.args.writer, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("WriteUInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := tt.args.writer.Bytes(); !bytes.Equal(gotWriter, tt.want) {
				t.Errorf("WriteUInt32() = %v, want %v", gotWriter, tt.want)
			}
		})
	}
}

func TestReadUInt32(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    uint32
		wantErr bool
	}{
		{name: "{0, 0, 0, 0} => 0", args: args{bytes.NewReader([]byte{0, 0, 0, 0})}, want: 0},
		{name: "{0, 0, 0, 1} => 1", args: args{bytes.NewReader([]byte{0, 0, 0, 1})}, want: 1},
		{name: "{0, 0, 1, 0} => 256", args: args{bytes.NewReader([]byte{0, 0, 1, 0})}, want: 256},
		{name: "{0, 0, 1, 0, 0} => 256", args: args{bytes.NewReader([]byte{0, 0, 1, 0, 0})}, want: 256},
		{name: "{1, 0, 0, 0} => 16777216", args: args{bytes.NewReader([]byte{1, 0, 0, 0})}, want: 16777216},
		{name: "{} => panic", args: args{bytes.NewReader([]byte{})}, want: 0, wantErr: true},
		{name: "{0, 1, 1} => err", args: args{bytes.NewReader([]byte{0, 1, 1})}, want: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadUInt32(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadUInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadUInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}
