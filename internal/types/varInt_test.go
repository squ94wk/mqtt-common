package types

import (
	"bytes"
	"io"
	"testing"
)

func TestWriteVarInt(t *testing.T) {
	type args struct {
		value uint32
	}
	tests := []struct {
		name       string
		args       args
		wantWriter []byte
		wantErr    bool
	}{
		{name: "0 => 0000 0000", args: args{0}, wantWriter: []byte{0}},
		{name: "1 => 0000 0001", args: args{1}, wantWriter: []byte{1}},
		{name: "127 => 0111 1111", args: args{127}, wantWriter: []byte{127}},
		{name: "128 => 1000 0000  0000 0001", args: args{128}, wantWriter: []byte{128, 1}},
		{name: "16,383 => 1111 1111  0111 1111", args: args{16383}, wantWriter: []byte{255, 127}},
		{name: "16,384 => 1000 0000  1000 0000  0000 0001", args: args{16384}, wantWriter: []byte{128, 128, 1}},
		{name: "2,097,151 => 1111 1111  1111 1111  0111 1111", args: args{2097151}, wantWriter: []byte{255, 255, 127}},
		{name: "2,097,152 => 1000 0000  1000 0000  1000 0000  0000 0001", args: args{2097152}, wantWriter: []byte{128, 128, 128, 1}},
		{name: "268,435,455 => 1111 1111  1111 1111  1111 1111  0111 1111", args: args{268435455}, wantWriter: []byte{255, 255, 255, 127}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := WriteVarInt(writer, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("WriteVarInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.Bytes(); !bytes.Equal(gotWriter, tt.wantWriter) {
				t.Errorf("WriteVarInt() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestReadVarInt(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    uint32
		wantErr bool
	}{
		{name: "0000 0000 => 0", args: args{bytes.NewReader([]byte{0})}, want: 0},
		{name: "0000 0001 => 1", args: args{bytes.NewReader([]byte{1})}, want: 1},
		{name: "0111 1111 => 127", args: args{bytes.NewReader([]byte{127})}, want: 127},
		{name: "1000 0000  0000 0001 => 128", args: args{bytes.NewReader([]byte{128, 1})}, want: 128},
		{name: "1111 1111  0111 1111 => 16,383", args: args{bytes.NewReader([]byte{255, 127})}, want: 16383},
		{name: "1000 0000  1000 0000  0000 0001 => 16,384", args: args{bytes.NewReader([]byte{128, 128, 1})}, want: 16384},
		{name: "1111 1111  1111 1111  0111 1111 => 2,097,151", args: args{bytes.NewReader([]byte{255, 255, 127})}, want: 2097151},
		{name: "1000 0000  1000 0000  1000 0000  0000 0001 => 2,097,152", args: args{bytes.NewReader([]byte{128, 128, 128, 1})}, want: 2097152},
		{name: "1111 1111  1111 1111  1111 1111  0111 1111 => 268,435,455", args: args{bytes.NewReader([]byte{255, 255, 255, 127})}, want: 268435455},
		{name: "5 bytes => err (exceeds max)", args: args{bytes.NewReader([]byte{128, 128, 128, 128, 1})}, want: 0, wantErr: true},
		{name: "1000 0000 => err (expecting next byte)", args: args{bytes.NewReader([]byte{128})}, want: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadVarInt(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadVarInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadVarInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
