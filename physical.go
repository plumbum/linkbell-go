package linkbell

import (
	"fmt"
	"io"
)

const (
	frameStart = 0x2B
	frameStop  = 0x4B

	preambule1 = 0xB2
	preambule2 = 0x4D
)

type PhysicalWriter struct {
	Preambule    []byte
	PreambuleLen int
	w            io.Writer
}

func NewPhysicalWriter(w io.Writer) *PhysicalWriter {
	enc := new(PhysicalWriter)
	enc.w = w
	enc.Preambule = []byte{preambule1, preambule2}
	enc.PreambuleLen = 4
	return enc
}

func (ce *PhysicalWriter) WritePreambule() (n int, err error) {
	n = 0
	for i := 0; i < ce.PreambuleLen; i++ {
		n1, err := ce.w.Write(ce.Preambule)
		n += n1
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (ce *PhysicalWriter) Write(p []byte) (n int, err error) {
	var n1 int
	n = 0

	// Write frame start byte
	n1, err = ce.w.Write([]byte{frameStart})
	n += n1
	if err != nil {
		return n, nil
	}

	// Write encoded packet
	for _, b := range p {
		hi, low := encodeNibbles(b)
		n1, err = ce.w.Write([]byte{hi, low})
		n += n1
		if err != nil {
			return n, err
		}
	}

	// Write frame stop byte
	n1, err = ce.w.Write([]byte{frameStop})
	n += n1
	return n, nil
}

type PhysicalReader struct {
	r              io.Reader
	inFrame        bool
	firstNibble    byte
	hasFirstNibble bool
}

func NewPhysicalReader(r io.Reader) *PhysicalReader {
	dec := new(PhysicalReader)
	dec.r = r
	return dec
}

func (cd *PhysicalReader) Reset() {
	cd.inFrame = false
	cd.hasFirstNibble = false
}

func (cd *PhysicalReader) InFrame() bool {
	return cd.inFrame
}

func (cd *PhysicalReader) Read(p []byte) (n int, err error) {
	n = 0
	buf := make([]byte, 1)
	for {
		if n >= len(p) {
			return n, nil
		}
		// Read one byte
		n1, err := cd.r.Read(buf)
		if err != nil {
			cd.Reset()
			return n, err
		}
		if n1 == 0 {
			return n, nil
		}
		b := buf[0]

		if !cd.inFrame { // Wait for frame start
			if b == frameStart {
				cd.inFrame = true
				cd.hasFirstNibble = false
			}
		} else {
			// Frame body
			if b == frameStop { // Frame end
				if cd.hasFirstNibble {
					cd.Reset()
					return n, fmt.Errorf("Unexpected frame end")
				}
				cd.inFrame = false
				return n, nil
			}
			if !cd.hasFirstNibble { // First nibble processing
				nibble, ok := decodeNibble(b)
				if !ok {
					cd.Reset()
					return n, fmt.Errorf("Undecodable byte 0x%02x", b)
				}
				cd.firstNibble = nibble
				cd.hasFirstNibble = true
			} else {
				nibble, ok := decodeNibble(b)
				if !ok {
					cd.Reset()
					return n, fmt.Errorf("Undecodable byte 0x%02x", b)
				}
				p[n] = joinNibbles(cd.firstNibble, nibble)
				n++
				cd.hasFirstNibble = false
			}
		}
	}

	return n, nil
}
