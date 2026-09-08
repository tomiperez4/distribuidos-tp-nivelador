package packet

import "errors"

type FinPacket struct{}

func NewFinPacket() Packet {
	return &FinPacket{}
}

func finFromBytes(payload []byte) (Packet, error) {
	if len(payload) > 0 {
		return nil, errors.New("unexpected payload in fin packet. Should be len == 0")
	}
	return NewFinPacket(), nil
}

func (fPkt *FinPacket) ToBytes() ([]byte, error) {
	return buildFrame(FIN_OPCODE, nil), nil
}
