package types

func StringSize(s string) uint32 {
	return 2 + uint32(len(s))
}

func BinarySize(b []byte) uint32 {
	return 2 + uint32(len(b))
}

func UInt16Size(i uint16) uint32 {
	return 2
}

func UInt32Size(i uint32) uint32 {
	return 4
}

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
