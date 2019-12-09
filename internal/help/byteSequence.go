package help

import (
	"fmt"

	"github.com/go-test/deep"
)

//CompareMode is an alias for the different modes byteSequences can be compared.
type CompareMode int

//Compare modes
const (
	InOrder CompareMode = iota
	Reversed
	AnyOrder
)

type segment struct {
	data []byte
	mode CompareMode
}
type sequence struct {
	data []ByteSequence
	mode CompareMode
}

//ByteSequence defines a sequence of segments that is compared to a byte slice.
type ByteSequence interface {
	Length() int
	Bytes() []byte
	Segments() []ByteSequence
	Mode() CompareMode
}

//Length returns the length of a segment in bytes.
func (s segment) Length() int {
	return len(s.data)
}

//Segments returns a slice of the segments in a segment.
func (s segment) Segments() []ByteSequence {
	return []ByteSequence{s}
}

//Bytes returns the segments concatenated as one byte slice.
func (s segment) Bytes() []byte {
	return s.data
}

//Mode returns the compare mode.
func (s segment) Mode() CompareMode {
	return s.mode
}

//Length returns the length of a sequence in bytes.
func (s sequence) Length() int {
	total := 0
	for _, segment := range s.data {
		total += segment.Length()
	}
	return total
}

//Segments returns the segments in a sequence.
func (s sequence) Segments() []ByteSequence {
	return s.data
}

//Bytes returns the sequence concatenated as one byte slice.
func (s sequence) Bytes() []byte {
	buf := make([]byte, 0, s.Length())
	for _, seg := range s.data {
		buf = append(buf, seg.Bytes()...)
	}
	return buf
}

//Mode returns the compare mode.
func (s sequence) Mode() CompareMode {
	return s.mode
}

//Match compares a byte slice against a ByteSequence.
func Match(seq ByteSequence, binary []byte) error {
	if len(binary) != seq.Length() {
		return fmt.Errorf("length wanted: %d, but got %d", len(binary), seq.Length())
	}

	if len(seq.Segments()) == 1 {
		if diff := deep.Equal(binary, seq.Bytes()); diff != nil {
			return fmt.Errorf("%v", diff)
		}
		return nil
	}

	offset := 0
	if seq.Mode() == AnyOrder {
		for i, seg := range seq.Segments() {
			if Match(seg, binary[:seg.Length()]) == nil {
				rest := append(seq.Segments()[:i], seq.Segments()[i+1:]...)

				if Match(NewByteSequence(AnyOrder, rest...), binary[seg.Length():]) == nil {
					return nil
				}
			}
		}
		return fmt.Errorf("sequence doesn't match bytes: attempted to match every combination of segments")
	}

	var successes []interface{}
	for _, p := range seq.Segments() {
		var next []byte
		switch seq.Mode() {
		case InOrder:
			next = binary[offset : offset+p.Length()]
			offset += p.Length()
		case Reversed:
			next = binary[len(binary)-offset-p.Length() : len(binary)-offset]
			offset += p.Length()
		default:
			panic("only modes 'IN_ORDER' , 'Reversed' are currently supported")
		}
		if diff := Match(p, next); diff != nil {
			return fmt.Errorf("match failed: %d matching: %v, sequence doesn't match bytes: %v", len(successes), successes, diff)
		}

		successes = append(successes, fmt.Sprintf("segment of length %d matches", p.Length()))
	}

	return nil
}

//NewByteSequence is a constructor for a ByteSequence concatenating a number of ByteSequences.
func NewByteSequence(mode CompareMode, segments ...ByteSequence) ByteSequence {
	return sequence{mode: mode, data: segments}
}

//NewByteSegment is a constructor for a ByteSequence concatenating a number of ByteSegments.
func NewByteSegment(segs ...[]byte) ByteSequence {
	return segment{data: Concat(segs...), mode: InOrder}
}

//Concat is an auxiliary function to concatenate byte slices.
func Concat(arrays ...[]byte) []byte {
	length := 0
	for _, a := range arrays {
		length += len(a)
	}

	buf := make([]byte, length)

	offset := 0
	for _, a := range arrays {
		copy(buf[offset:offset+len(a)], a)
		offset += len(a)
	}

	return buf
}
