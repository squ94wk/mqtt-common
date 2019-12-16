package help

import (
	"fmt"

	"github.com/go-test/deep"
)

type CompareMode int

const (
	IN_ORDER CompareMode = iota
	REVERSED
	ANY_ORDER
)

type segment struct {
	data []byte
	mode CompareMode
}
type sequence struct {
	data []ByteSequence
	mode CompareMode
}

type ByteSequence interface {
	Length() int
	Bytes() []byte
	Segments() []ByteSequence
	Mode() CompareMode
}

func (s segment) Length() int {
	return len(s.data)
}

func (s segment) Segments() []ByteSequence {
	return []ByteSequence{s}
}

func (s segment) Bytes() []byte {
	return s.data
}

func (s segment) Mode() CompareMode {
	return s.mode
}

func (s sequence) Length() int {
	total := 0
	for _, segment := range s.data {
		total += segment.Length()
	}
	return total
}

func (s sequence) Segments() []ByteSequence {
	return s.data
}

func (s sequence) Bytes() []byte {
	buf := make([]byte, 0, s.Length())
	for _, seg := range s.data {
		buf = append(buf, seg.Bytes()...)
	}
	return buf
}

func (s sequence) Mode() CompareMode {
	return s.mode
}

type byteSequence struct {
	segments []ByteSequence
	mode     CompareMode
}

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
	if seq.Mode() == ANY_ORDER {
		for i, seg := range seq.Segments() {
			if Match(seg, binary[:seg.Length()]) == nil {
				rest := append(seq.Segments()[:i], seq.Segments()[i+1:]...)

				if Match(NewByteSequence(ANY_ORDER, rest...), binary[seg.Length():]) == nil {
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
		case IN_ORDER:
			next = binary[offset : offset+p.Length()]
			offset += p.Length()
		case REVERSED:
			next = binary[len(binary)-offset-p.Length() : len(binary)-offset]
			offset += p.Length()
		default:
			panic("only modes 'IN_ORDER' , 'REVERSED' are currently supported")
		}
		if diff := Match(p, next); diff != nil {
			return fmt.Errorf("match failed: %d matching: %v, sequence doesn't match bytes: %v", len(successes), successes, diff)
		}

		successes = append(successes, fmt.Sprintf("segment of length %d matches", p.Length()))
	}

	return nil
}

func NewByteSequence(mode CompareMode, segments ...ByteSequence) ByteSequence {
	return sequence{mode: mode, data: segments}
}

func NewByteSegment(segs ...[]byte) ByteSequence {
	return segment{data: Concat(segs...), mode: IN_ORDER}
}

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
