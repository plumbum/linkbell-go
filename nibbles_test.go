package linkbell

import (
	"bytes"
	"testing"
)

var (
	codeTable = []byte{
		0x55, 0x56, 0x59, 0x5A,
		0x65, 0x66, 0x69, 0x6A,
		0x95, 0x96, 0x99, 0x9A,
		0xA5, 0xA6, 0xA9, 0xAA,
	}
)

func TestEncodeNibble(t *testing.T) {
	for n := 0; n < 16; n++ {
		enc := encodeNibble(byte(n))
		if enc != codeTable[n] {
			t.Fatalf("encodeNibble: for nibble %d expected 0x%02x, got 0x%02x", n, codeTable[n], enc)
		}
	}
}

func TestDecodeNibble(t *testing.T) {
	for b := 0; b < 256; b++ {
		i := bytes.IndexByte(codeTable, byte(b))
		if i >= 0 {
			if n, ok := decodeNibble(byte(b)); ok {
				if int(n) != i {
					t.Fatalf("encodeNibble: for byte 0x%02x expected %d, got %d", b, i, n)
				}
			} else {
				t.Fatalf("encodeNibble: for byte 0x%02x expected %d, got error", b, i)
			}
		} else {
			if n, ok := decodeNibble(byte(b)); ok {
				t.Fatalf("encodeNibble: for byte 0x%02x expected error, got %d", n)
			}
		}
	}
}
