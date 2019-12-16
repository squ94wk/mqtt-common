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
	if i > 2<<21 {
		return 4
	}
	if i > 2<<14 {
		return 3
	}
	if i > 2<<7 {
		return 2
	}
	return 1
}
