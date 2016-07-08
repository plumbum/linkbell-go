package linkbell

import (
	"bytes"
	"testing"
)

var (
	channelPacket = []byte{0xE2, 0x5D, 0x06, 0x00, 0x00, 0x64, 0x03, 0xE8, 0xEF, 0xBE, 0xAD, 0xDE}
)

func TestChannelWriter_Write(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})

	c := NewChannelWriter(buf)
	c.NetworkId = 10
	c.ReceiverId = 100
	c.SenderId = 1000

	// Write frame
	n, err := c.Write(rawPacket)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buf.Bytes(), channelPacket) {
		t.Fatalf("ChannelWriter: expected %v, got %v (%d)", channelPacket, buf.Bytes(), n)
	}
}

func TestChannelReader_Read(t *testing.T) {

	buf := bytes.NewBuffer([]byte{})

	cw := NewChannelWriter(buf)
	cw.NetworkId = 10
	cw.ReceiverId = 100
	cw.SenderId = 1000
	cw.Flags = 0x5A
	cw.Repeaters = []DeviceId{555, 666, 777}

	// Write frame
	n, err := cw.Write(rawPacket)
	if err != nil {
		t.Fatal(err)
	}

	outBuf := make([]byte, 256)

	cr := NewChannelReader(buf)
	cr.NetworkId = 10
	n, err = cr.Read(outBuf)
	if err != nil {
		t.Fatal(err)
	}

	if cr.ReceiverId != cw.ReceiverId {
		t.Errorf("ChannelReader: expected ReceiverId %d, got %d", cw.ReceiverId, cr.ReceiverId)
	}
	if cr.SenderId != cw.SenderId {
		t.Errorf("ChannelReader: expected SenderId %d, got %d", cw.SenderId, cr.SenderId)
	}

	if cr.Flags != cw.Flags {
		t.Errorf("ChannelReader: expected Flags %02x, got %02x", cw.Flags, cr.Flags)
	}

	if len(cr.Repeaters) != len(cw.Repeaters) {
		t.Errorf("ChannelReader: expected %d repeaters, got %d repeaters", len(cw.Repeaters), len(cr.Repeaters))
	} else {
		for i, r := range cr.Repeaters {
			if r != cw.Repeaters[i] {
				t.Errorf("ChannelReader: expected Repeater[%d] %d, got %d", i, cw.Repeaters[i], r)
			}
		}
	}

	if !bytes.Equal(outBuf[:n], rawPacket) {
		t.Errorf("ChannelReader: expected data %v (%d), got %v (%d)", rawPacket, len(rawPacket), outBuf[:n], n)
	}
}
