package types

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestWriteBinary(t *testing.T) {
	longPayload := make([]byte, ^uint16(0))
	type args struct {
		writer bytes.Buffer
		value  []byte
	}
	tests := []struct {
		name       string
		args       args
		wantWriter []byte
		wantErr    bool
	}{
		{name: "{} => {0, 0}", args: args{bytes.Buffer{}, []byte{}}, wantWriter: []byte{0, 0}},
		{name: "{127} => {0, 1, 127}", args: args{bytes.Buffer{}, []byte{127}}, wantWriter: []byte{0, 1, 127}},
		{name: "{1, 2, 3} => {0, 3, 1, 2, 3}", args: args{bytes.Buffer{}, []byte{1, 2, 3}}, wantWriter: []byte{0, 3, 1, 2, 3}},
		{name: "{longPayload...} => {255, 255, 000000...}", args: args{bytes.Buffer{}, longPayload}, wantWriter: append([]byte{255, 255}, longPayload...)},
		{name: "{tooLongPayload...} => {255, 255}", args: args{bytes.Buffer{}, make([]byte, 1<<16)}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteBinary(&tt.args.writer, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("WriteBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := tt.args.writer.Bytes(); !bytes.Equal(gotWriter, tt.wantWriter) {
				t.Errorf("WriteBinary() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestReadBinary(t *testing.T) {
	longPayload := make([]byte, ^uint16(0))
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{name: "{0, 0} => {}", args: args{bytes.NewReader([]byte{0, 0})}, want: []byte{}},
		{name: "{0, 1, 127} => {127}", args: args{bytes.NewReader([]byte{0, 1, 127})}, want: []byte{127}},
		{name: "{0, 3, 1, 2, 3} => {1, 2, 3}", args: args{bytes.NewReader([]byte{0, 3, 1, 2, 3})}, want: []byte{1, 2, 3}},
		{name: "{0, 2, 1, 2, 3} => {1, 2}", args: args{bytes.NewReader([]byte{0, 2, 1, 2, 3})}, want: []byte{1, 2}},
		{name: "{255, 255, 000000...} => {longPayload...}", args: args{bytes.NewReader(append([]byte{255, 255}, longPayload...))}, want: longPayload},
		{name: "{0, 4, 1, 2, 3} => err", args: args{bytes.NewReader([]byte{0, 4, 1, 2, 3})}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadBinary(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}
