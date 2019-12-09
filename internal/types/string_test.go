package types

import (
	"bytes"
	"io"
	"testing"
)

func TestReadString(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "{0, 0} => empty", args: args{bytes.NewReader([]byte{0, 0})}, want: ""},
		{name: "{0, 1, 'a'} => 'a'", args: args{bytes.NewReader([]byte{0, 1, 'a'})}, want: "a"},
		{name: "{0, 12, 'longerString'...} => 'longerString'", args: args{bytes.NewReader(append([]byte{0, 12}, []byte("longerString")...))}, want: "longerString"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadString(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteString(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{name: "empty string => {0, 0}", args: args{""}, want: []byte{0, 0}},
		{name: "'a' => {0, 1, 'a'}", args: args{"a"}, want: []byte{0, 1, 'a'}},
		{name: "'longerString' => {0, 12, 'longString'...}", args: args{"longerString"}, want: append([]byte{0, 12}, []byte("longerString")...)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if _, err := WriteStringTo(writer, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("WriteStringTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.Bytes(); !bytes.Equal(gotWriter, tt.want) {
				t.Errorf("WriteStringTo() = %v, want %v", gotWriter, tt.want)
			}
		})
	}
}
