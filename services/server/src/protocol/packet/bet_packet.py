from .packet import Packet, build_frame
from lottery import Bet
from .bet_codec import encode_bet, decode_bet
from .opcode import Opcode


class BetPacket(Packet):
    def __init__(self, bets: list[Bet]):
        self.bets = bets

    def to_bytes(self) -> bytes:
        payload = b"".join(encode_bet(bet) for bet in self.bets)
        return build_frame(Opcode.BET, payload)

def bet_pkt_from_bytes(data: bytes, agency_id: int) -> Packet:
    offset = 0
    bets = []

    while offset < len(data):
        bet, consumed = decode_bet(data[offset:], agency_id)
        bets.append(bet)
        offset += consumed

    return BetPacket(bets)
