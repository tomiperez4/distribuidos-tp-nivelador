package packet

type OpCode uint8

// No existen enums en GO. Aplico esto que es "similar"
const (
	BET_OPCODE OpCode = iota
	FIN_OPCODE
	ACK_OPCODE
	HELLO_OPCODE
)
