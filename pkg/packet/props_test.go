package packet

import (
	"bytes"
	"io"
	"testing"

	"github.com/go-test/deep"
	"github.com/squ94wk/mqtt-common/internal/help"
)

func TestReadProperties(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    Properties
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
			want: NewProperties(
				Property{propID: SessionExpiryInterval, payload: Int32PropPayload(17)},
				Property{propID: AssignedClientIdentifier, payload: StringPropPayload("client")},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readProperties(tt.args.reader)
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

//TODO: TestWritePropsTo
