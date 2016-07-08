package linkbell

func joinNibbles(hi, low byte) byte {
	return (hi << 4) | (low & 0x0f)
}

func encodeNibbles(b byte) (hi, low byte) {
	return encodeNibble((b >> 4) & 0x0F),
		encodeNibble((b >> 0) & 0x0F)
}

/*
	codeTable = []byte{
		0x55, 0x56, 0x59, 0x5A,
		0x65, 0x66, 0x69, 0x6A,
		0x95, 0x96, 0x99, 0x9A,
		0xA5, 0xA6, 0xA9, 0xAA,
	}
*/
func encodeNibble(n byte) (b byte) {
	if n&1 == 1 {
		b |= 0x02
	} else {
		b |= 0x01
	}
	if n&2 == 2 {
		b |= 0x08
	} else {
		b |= 0x04
	}
	if n&4 == 4 {
		b |= 0x20
	} else {
		b |= 0x10
	}
	if n&8 == 8 {
		b |= 0x80
	} else {
		b |= 0x40
	}
	return b
}

func decodeNibble(b byte) (n byte, ok bool) {
	goodPairs := 0
	// bit 0
	if b&0x03 == 0x01 {
		goodPairs++
	}
	if b&0x03 == 0x02 {
		goodPairs++
		n |= 1
	}
	// bit 1
	if b&0x0C == 0x04 {
		goodPairs++
	}
	if b&0x0C == 0x08 {
		goodPairs++
		n |= 2
	}
	// bit 2
	if b&0x30 == 0x10 {
		goodPairs++
	}
	if b&0x30 == 0x20 {
		goodPairs++
		n |= 4
	}
	// bit 3
	if b&0xC0 == 0x40 {
		goodPairs++
	}
	if b&0xC0 == 0x80 {
		goodPairs++
		n |= 8
	}
	if goodPairs != 4 {
		return n, false
	}
	return n, true
}
