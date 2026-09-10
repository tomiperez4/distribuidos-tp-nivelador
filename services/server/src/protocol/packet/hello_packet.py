from .packet import Packet, build_frame, AGENCY_ID_SIZE
from .opcode import Opcode


class HelloPacket(Packet):
    def __init__(self, agency_id: int):
        self.agency_id = agency_id

    def to_bytes(self) -> bytes:
        return build_frame(Opcode.HELLO, self.agency_id.to_bytes(AGENCY_ID_SIZE, byteorder='big'))

def hello_pkt_from_bytes(data: bytes) -> HelloPacket:
    if len(data) != AGENCY_ID_SIZE:
        raise Exception(f"hello packet: expected {AGENCY_ID_SIZE} byte(s) but got {len(data)}")
    return HelloPacket(data[0])
