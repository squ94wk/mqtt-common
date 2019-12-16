package packet

import (
	"bytes"
	"io"
	"testing"

	"github.com/go-test/deep"
	"github.com/squ94wk/mqtt-common/internal/help"
)

func NewProps(props ...Property) map[PropId][]Property {
	properties := make(map[PropId][]Property)
	for _, p := range props {
		if withId, ok := properties[p.PropId()]; ok {
			properties[p.PropId()] = append(withId, p)
		} else {
			properties[p.PropId()] = []Property{p}
		}
	}

	return properties
}

func TestReadProperties(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    map[PropId][]Property
		wantErr bool
	}{
		{
			name: "",
			args: args{
				reader: bytes.NewReader(help.Concat(
					[]byte{14},
					[]byte{byte(SessionExpiryInterval), 0, 0, 0, 17},
					[]byte{byte(AssignedClientIdentifier), 0, 6, 'c', 'l', 'i', 'e', 'n', 't'},
				)),
			},
			want: NewProps(
				NewInt32Prop(SessionExpiryInterval, 10),
				NewStringProp(AssignedClientIdentifier, "client"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make(map[PropId][]Property)
			err := readProperties(tt.args.reader, got)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ReadProperties() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Error(diff)
			}
		})
	}
}
