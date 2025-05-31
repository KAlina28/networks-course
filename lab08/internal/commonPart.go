package internal

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
)

type Packet struct {
	SeqNum   uint32
	Checksum uint32
	Data     []byte
}

func MakePacket(seq uint32, data []byte) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, seq)
	checksum := CalculateChecksum(seq, data)
	binary.Write(buf, binary.BigEndian, checksum)
	buf.Write(data)
	return buf.Bytes()
}

func MakeACK(seq uint32) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, seq)
	binary.Write(buf, binary.BigEndian, CalculateChecksum(seq, nil))
	return buf.Bytes()
}

func CalculateChecksum(seq uint32, data []byte) uint32 {
	h := md5.New()
	binary.Write(h, binary.BigEndian, seq)
	h.Write(data)
	var sum uint32
	binary.Read(bytes.NewReader(h.Sum(nil)), binary.BigEndian, &sum)
	return sum
}

func ParsePacket(data []byte) (*Packet, error) {
	if len(data) < 8 {
		return nil, errors.New("packet too short")
	}
	return &Packet{
		SeqNum:   binary.BigEndian.Uint32(data[:4]),
		Checksum: binary.BigEndian.Uint32(data[4:8]),
		Data:     data[8:],
	}, nil
}

func VerifyChecksum(pkt *Packet) bool {
	return pkt.Checksum == CalculateChecksum(pkt.SeqNum, pkt.Data)
}
