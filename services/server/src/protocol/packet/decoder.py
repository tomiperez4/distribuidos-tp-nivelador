from . import Packet, Opcode, bet_pkt_from_bytes, fin_pkt_from_bytes, ack_pkt_from_bytes


def from_bytes(opcode: int, buff: bytes) -> Packet:
    if opcode == Opcode.BET:
        return bet_pkt_from_bytes(buff)
    elif opcode == Opcode.FIN:
        return fin_pkt_from_bytes(buff)
    elif opcode == Opcode.ACK:
        return ack_pkt_from_bytes(buff)
    else:
        raise Exception(f"opcode {opcode} not supported")