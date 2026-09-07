from packet import Opcode
from packet.packet import Packet, build_frame


class FinPacket(Packet):
    def to_bytes(self) -> bytes:
        return build_frame(Opcode.FIN, b"")

def fin_pkt_from_bytes(data: bytes) -> FinPacket:
    if len(data) > 0:
        raise Exception(f"fin packet: expected no data but got {len(data)}")
    return FinPacket()