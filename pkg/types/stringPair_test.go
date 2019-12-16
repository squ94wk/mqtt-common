package types

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

//func TestEncodeStringPair(t *testing.T) {
//type args struct {
//pair StringPair
//}
//tests := []struct {
//name string
//args args
//want []byte
//}{
//{name: "'a', 'b' => {0, 1, 'a', 0, 1, 'b'}", args: args{StringPair{"a", "b"}}, want: []byte{0, 1, 'a', 0, 1, 'b'}},
//{name: "'', '' => {0, 0, 0, 0}", args: args{StringPair{"", ""}}, want: []byte{0, 0, 0, 0}},
//{name: "'abc', '' => {0, 3, 'a', 'b', 'c', 0, 0}", args: args{StringPair{"abc", ""}}, want: []byte{0, 3, 'a', 'b', 'c', 0, 0}},
//}
//for _, tt := range tests {
//t.Run(tt.name, func(t *testing.T) {
//if got := EncodeStringPair(tt.args.pair); !reflect.DeepEqual(got, tt.want) {
//t.Errorf("EncodeStringPair() = %v, want %v", got, tt.want)
//}
//})
//}
//}

//func TestDecodeStringPair(t *testing.T) {
//type args struct {
//buf []byte
//}
//tests := []struct {
//name    string
//args    args
//want    StringPair
//wantErr bool
//}{
//{name: "{0, 1, 'a', 0, 1, 'b'} => 'a', 'b'", args: args{[]byte{0, 1, 'a', 0, 1, 'b'}}, want: StringPair{"a", "b"}},
//{name: "{0, 0, 0, 0} => '', ''", args: args{[]byte{0, 0, 0, 0}}, want: StringPair{"", ""}},
//{name: "{0, 3, 'a', 'b', 'c', 0, 0} => 'abc', ''", args: args{[]byte{0, 3, 'a', 'b', 'c', 0, 0}}, want: StringPair{"abc", ""}},
//{name: "{0, 1, 'a', 'b', 'c', 0, 1, 'd'} => err", args: args{[]byte{0, 1, 'a', 'b', 'c', 0, 1, 'd'}}, wantErr: true},
//}
//for _, tt := range tests {
//t.Run(tt.name, func(t *testing.T) {
//got, err := DecodeStringPair(tt.args.buf)
//if (err != nil) != tt.wantErr {
//t.Errorf("DecodeStringPair() error = %v, wantErr %v", err, tt.wantErr)
//return
//}
//if !reflect.DeepEqual(got, tt.want) {
//t.Errorf("DecodeStringPair() = %v, want %v", got, tt.want)
//}
//})
//}
//}

func TestWriteStringPair(t *testing.T) {
	type args struct {
		value StringPair
	}
	tests := []struct {
		name       string
		args       args
		wantWriter []byte
		wantErr    bool
	}{
		{name: "'a', 'b' => {0, 1, 'a', 0, 1, 'b'}", args: args{StringPair{"a", "b"}}, wantWriter: []byte{0, 1, 'a', 0, 1, 'b'}},
		{name: "'', '' => {0, 0, 0, 0}", args: args{StringPair{"", ""}}, wantWriter: []byte{0, 0, 0, 0}},
		{name: "'abc', '' => {0, 3, 'a', 'b', 'c', 0, 0}", args: args{StringPair{"abc", ""}}, wantWriter: []byte{0, 3, 'a', 'b', 'c', 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := WriteStringPair(writer, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("WriteStringPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.Bytes(); !bytes.Equal(gotWriter, tt.wantWriter) {
				t.Errorf("WriteStringPair() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestReadStringPair(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    StringPair
		wantErr bool
	}{
		{name: "{0, 1, 'a', 0, 1, 'b'} => 'a', 'b'", args: args{bytes.NewReader([]byte{0, 1, 'a', 0, 1, 'b'})}, want: StringPair{"a", "b"}},
		{name: "{0, 0, 0, 0} => '', ''", args: args{bytes.NewReader([]byte{0, 0, 0, 0})}, want: StringPair{"", ""}},
		{name: "{0, 3, 'a', 'b', 'c', 0, 0} => 'abc', ''", args: args{bytes.NewReader([]byte{0, 3, 'a', 'b', 'c', 0, 0})}, want: StringPair{"abc", ""}},
		{name: "{0, 1, 'a', 'b', 'c', 0, 1, 'd'} => err", args: args{bytes.NewReader([]byte{0, 1, 'a', 'b', 'c', 0, 1, 'd'})}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadStringPair(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadStringPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadStringPair() = %v, want %v", got, tt.want)
			}
		})
	}
}
