package packet

import (
	"encoding/binary"
	"errors"
)

const (
	HEADER_LENGTH = 5
)

type Packet interface {
	// ToBytes transform the packet into a slice of bytes. Returns also an error type, in case it occurs
	ToBytes() ([]byte, error)
}

func buildFrame(opCode OpCode, payload []byte) []byte {
	buff := make([]byte, 0, HEADER_LENGTH+len(payload))

	buff = append(buff, byte(opCode))
	buff = binary.BigEndian.AppendUint32(buff, uint32(len(payload)))
	buff = append(buff, payload...)

	return buff
}

func ParseHeader(buff []byte) (OpCode, uint32, error) {
	if len(buff) < HEADER_LENGTH {
		return 0, 0, errors.New("packet header too short")
	}
	opCode := OpCode(buff[0])
	payloadLen := binary.BigEndian.Uint32(buff[1:HEADER_LENGTH])

	return opCode, payloadLen, nil
}

func FromBytes(opCode OpCode, payload []byte) (Packet, error) {
	switch opCode {
	case BET_OPCODE:
		return betFromBytes(payload)
	case FIN_OPCODE:
		return finFromBytes(payload)
	case ACK_OPCODE:
		return ackFromBytes(payload)
	default:
		return nil, errors.New("invalid packet opcode")
	}
}
