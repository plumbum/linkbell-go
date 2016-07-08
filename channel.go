package linkbell

import (
	"bytes"
	"errors"
	"io"
)

const (
	WireId      NetworkId = 0x0000
	EmptyId     DeviceId  = 0x0000
	BroadcastId DeviceId  = 0x7FFF
)

type NetworkId int16

type DeviceId int16

type Channel struct {
	NetworkId  NetworkId
	Flags      byte
	ReceiverId DeviceId
	SenderId   DeviceId
	Repeaters  []DeviceId
}

type ChannelWriter struct {
	Channel
	w io.Writer
}

func NewChannelWriter(w io.Writer) *ChannelWriter {
	c := new(ChannelWriter)
	c.w = w
	return c
}

func (c *ChannelWriter) Write(p []byte) (n int, err error) {
	var hi, low byte

	buf := bytes.NewBuffer([]byte{})

	headerLen := 6 + len(c.Repeaters)*2
	buf.WriteByte(byte(headerLen))

	buf.WriteByte(c.Flags)

	hi, low = word2byte(uint16(c.ReceiverId))
	buf.WriteByte(hi)
	buf.WriteByte(low)

	hi, low = word2byte(uint16(c.SenderId))
	buf.WriteByte(hi)
	buf.WriteByte(low)

	for _, r := range c.Repeaters {
		hi, low = word2byte(uint16(r))
		buf.WriteByte(hi)
		buf.WriteByte(low)
	}
	buf.Write(p)

	crc := new(crc)
	crc.init(uint16(c.NetworkId))
	crc.pushBytes(buf.Bytes())

	n, err = c.w.Write(append([]byte{crc.low, crc.high}, buf.Bytes()...))
	return n, err
}

func word2byte(w uint16) (hi, low byte) {
	return byte((w >> 8) & 0xFF), byte(w & 0xFF)
}

type ChannelReader struct {
	Channel
	r io.Reader
}

func NewChannelReader(r io.Reader) *ChannelReader {
	c := new(ChannelReader)
	c.r = r
	return c
}

func (c *ChannelReader) Reset() {
	c.ReceiverId = 0
	c.SenderId = 0
}

func (c *ChannelReader) Read(p []byte) (n int, err error) {
	buf := make([]byte, 256)
	n1, err := c.r.Read(buf)
	if err != nil {
		c.Reset()
		return 0, err
	}
	data := buf[2:n1]

	crc := new(crc)
	crc.init(uint16(c.NetworkId))
	crc.pushBytes(data)
	if buf[0] != crc.low || buf[1] != crc.high {
		return 0, errors.New("CRC error")
	}

	headerLen := int(data[0])
	if len(data) < headerLen {
		return 0, errors.New("Packet too small")
	}
	dataLen := len(data) - headerLen
	if len(p) < dataLen {
		return 0, errors.New("Data too long")
	}

	c.Flags = data[1]
	c.ReceiverId = DeviceId(uint16(data[2])<<8 | uint16(data[3]))
	c.SenderId = DeviceId(uint16(data[4])<<8 | uint16(data[5]))

	repLen := (headerLen - 6) / 2
	c.Repeaters = make([]DeviceId, repLen)
	for i := 0; i < repLen; i++ {
		idx := i*2 + 6
		c.Repeaters[i] = DeviceId(uint16(data[idx])<<8 | uint16(data[idx+1]))
	}
	copy(p, data[headerLen:])

	return dataLen, nil
}
