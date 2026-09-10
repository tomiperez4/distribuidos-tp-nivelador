package packet

type HelloPacket struct {
	AgencyId uint8
}

func NewHelloPacket(agencyId uint8) Packet {
	return &HelloPacket{AgencyId: agencyId}
}

func (hPkt *HelloPacket) ToBytes() ([]byte, error) {
	payload := make([]byte, AGENCY_ID_SIZE)
	payload[0] = hPkt.AgencyId
	return buildFrame(HELLO_OPCODE, payload), nil
}
