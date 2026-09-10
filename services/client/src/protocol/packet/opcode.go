package packet

type OpCode uint8

const (
	BET_OPCODE OpCode = iota
	FIN_OPCODE
	ACK_OPCODE
	HELLO_OPCODE
)
