from . import Packet, Opcode, bet_pkt_from_bytes, fin_pkt_from_bytes, ack_pkt_from_bytes, hello_pkt_from_bytes


def from_bytes(opcode: int, buff: bytes, agency_id: int | None = None) -> Packet:
    if opcode == Opcode.BET:
        return bet_pkt_from_bytes(buff, agency_id)
    elif opcode == Opcode.FIN:
        return fin_pkt_from_bytes(buff)
    elif opcode == Opcode.ACK:
        return ack_pkt_from_bytes(buff)
    elif opcode == Opcode.HELLO:
        return hello_pkt_from_bytes(buff)
    else:
        raise Exception(f"opcode {opcode} not supported")