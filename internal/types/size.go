package types

//StringSize returns the length taken by a encoded string.
func StringSize(s string) uint32 {
	return 2 + uint32(len(s))
}

//BinarySize returns the length taken by a encoded byte array.
func BinarySize(b []byte) uint32 {
	return 2 + uint32(len(b))
}

const (
	//UInt32Size returns the length taken by a encoded 32 bit integer.
	UInt32Size uint32 = 4
	//UInt16Size returns the length taken by a encoded 16 bit integer.
	UInt16Size uint32 = 2
)

//VarIntSize returns the length taken by a encoded variable length integer.
func VarIntSize(i uint32) uint32 {
	switch {
	case i > 1<<22:
		return 4
	case i > 1<<15:
		return 3
	case i > 1<<8:
		return 2
	default:
		return 1
	}
}
