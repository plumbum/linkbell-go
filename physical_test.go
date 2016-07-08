package linkbell

import (
	"bytes"
	"testing"
)

var (
	rawPacket     = []byte{0xEF, 0xBE, 0xAD, 0xDE} // DEADBEEF
	encodedPacket = []byte{
		preambule1, preambule2,
		preambule1, preambule2,
		preambule1, preambule2,
		preambule1, preambule2,
		frameStart, 0xA9, 0xAA, 0x9A, 0xA9, 0x99, 0xA6, 0xA6, 0xA9, frameStop,
	}
)

func TestPhysicalWriter_Write(t *testing.T) {

	buf := bytes.NewBuffer([]byte{})

	enc := NewPhysicalWriter(buf)

	// Write preambule
	n, err := enc.WritePreambule()
	if err != nil {
		t.Fatal(err)
	}

	// Write frame
	n, err = enc.Write(rawPacket)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buf.Bytes(), encodedPacket) {
		t.Fatalf("PhysicalWriter: expected %v, got %v (%d)", encodedPacket, buf.Bytes(), n)
	}

}

func TestPhysicalReader_Read_2CorrectPacket(t *testing.T) {

	inBuf := bytes.NewBuffer(append(encodedPacket, encodedPacket...))
	dec := NewPhysicalReader(inBuf)

	outBuf := make([]byte, 8)

	n, err := dec.Read(outBuf)
	if err != nil {
		t.Fatal(err)
	}
	if n != 4 {
		t.Fatalf("PhysicalReader: expected 4 bytes, got %d", n)
	}
	if !bytes.Equal(rawPacket, outBuf[0:4]) {
		t.Fatalf("PhysicalReader: expected %v, got %v", rawPacket, outBuf)
	}

	n, err = dec.Read(outBuf)
	if err != nil {
		t.Fatal(err)
	}
	if n != 4 {
		t.Fatalf("PhysicalReader: expected 4 bytes, got %d", n)
	}
	if !bytes.Equal(rawPacket, outBuf[0:4]) {
		t.Fatalf("PhysicalReader: expected %v, got %v", rawPacket, outBuf)
	}
}

func TestPhysicalReader_Read_CorrectPacketInTwoStages(t *testing.T) {

	inBuf := bytes.NewBuffer(encodedPacket)
	dec := NewPhysicalReader(inBuf)

	outBuf := make([]byte, 2)

	n, err := dec.Read(outBuf)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("PhysicalReader: expected 2 bytes, got %d", n)
	}
	if !bytes.Equal(rawPacket[0:2], outBuf) {
		t.Fatalf("PhysicalReader: expected %v, got %v", rawPacket, outBuf)
	}

	n, err = dec.Read(outBuf)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("PhysicalReader: expected 2 bytes, got %d", n)
	}
	if !bytes.Equal(rawPacket[2:4], outBuf) {
		t.Fatalf("PhysicalReader: expected %v, got %v", rawPacket, outBuf)
	}
}

func TestPhysicalReader_Read_ErrorUndecodable(t *testing.T) {
	errorPacket := []byte{
		preambule1, preambule2,
		preambule1, preambule2,
		preambule1, preambule2,
		preambule1, preambule2,
		frameStart, 0xA9, 0xAA, 0x9A, 0xA9, 0x99, 0xA6, 0xA6, 0xA7, frameStop,
	}

	inBuf := bytes.NewBuffer(errorPacket)
	dec := NewPhysicalReader(inBuf)

	outBuf := make([]byte, 8)

	n, err := dec.Read(outBuf)
	if err == nil {
		t.Fatalf("PhysicalReader: expected error, but got packet (%d) %v", n, outBuf)
	}

	t.Logf("PhysicalReader got error: %v", err)
}

func TestPhysicalReader_Read_ErrorUnexpected(t *testing.T) {
	errorPacket := []byte{
		preambule1, preambule2,
		preambule1, preambule2,
		preambule1, preambule2,
		preambule1, preambule2,
		frameStart, 0xA9, 0xAA, 0x9A, 0xA9, 0x99, 0xA6, 0xA6, frameStop,
	}

	inBuf := bytes.NewBuffer(errorPacket)
	dec := NewPhysicalReader(inBuf)

	outBuf := make([]byte, 8)

	n, err := dec.Read(outBuf)
	if err == nil {
		t.Fatalf("PhysicalReader: expected error, but got packet (%d) %v", n, outBuf)
	}

	t.Logf("PhysicalReader got error: %v", err)
}
