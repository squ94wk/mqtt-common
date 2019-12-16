package help

import "testing"

func TestMatch(t *testing.T) {
	type args struct {
		seq    ByteSequence
		binary []byte
	}
	tests := []struct {
		name string
		args args
		diff bool
	}{
		{name: "one segment => true", args: args{binary: []byte{0, 1}, seq: NewByteSegment([]byte{0, 1})}, diff: true},
		{name: "one segment wrong => false", args: args{binary: []byte{1, 0}, seq: NewByteSegment([]byte{0, 1})}, diff: false},
		{name: "two segments => true", args: args{binary: []byte{0, 1, 2, 3}, seq: NewByteSequence(IN_ORDER, NewByteSegment([]byte{0, 1}), NewByteSegment([]byte{2, 3}))}, diff: true},
		{name: "two segments wrong => false", args: args{binary: []byte{2, 3, 0, 1}, seq: NewByteSequence(IN_ORDER, NewByteSegment([]byte{0, 1}), NewByteSegment([]byte{2, 3}))}, diff: false},
		{name: "reversed 2 => true", args: args{binary: []byte{2, 3, 0, 1}, seq: NewByteSequence(REVERSED, NewByteSegment([]byte{0, 1}), NewByteSegment([]byte{2, 3}))}, diff: true},
		{name: "any order => true", args: args{binary: []byte{2, 3, 0, 1, 4, 5, 7}, seq: NewByteSequence(ANY_ORDER, NewByteSegment([]byte{0, 1}), NewByteSegment([]byte{2, 3}), NewByteSegment([]byte{4, 5}), NewByteSegment([]byte{7}))}, diff: true},
		{name: "any order with duplicates => true", args: args{binary: []byte{2, 3, 0, 1, 0, 1}, seq: NewByteSequence(ANY_ORDER, NewByteSegment([]byte{0, 1}), NewByteSegment([]byte{2, 3}), NewByteSegment([]byte{0, 1}))}, diff: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := Match(tt.args.seq, tt.args.binary); (diff == nil) != tt.diff {
				t.Errorf("Match() = %v, want %v", diff, tt.diff)
			}
		})
	}
}
