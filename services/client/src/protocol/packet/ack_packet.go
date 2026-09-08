package packet

import "errors"

type AckPacket struct{}

func NewAckPacket() Packet {
	return &AckPacket{}
}

func ackFromBytes(payload []byte) (Packet, error) {
	if len(payload) > 0 {
		return nil, errors.New("unexpected payload in ack packet. Should be len == 0")
	}
	return NewAckPacket(), nil
}

func (aPkt *AckPacket) ToBytes() ([]byte, error) {
	return buildFrame(ACK_OPCODE, nil), nil
}
