from .opcode import Opcode
from .packet import Packet, build_frame


class AckPacket(Packet):
    def to_bytes(self) -> bytes:
        return build_frame(Opcode.ACK, b"")

def ack_pkt_from_bytes(data: bytes) -> AckPacket:
    if len(data) > 0:
        raise Exception(f"ack packet: expected no data but got {len(data)}")
    return AckPacket()
