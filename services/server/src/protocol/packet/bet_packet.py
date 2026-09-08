from .packet import Packet, build_frame
from lottery import Bet
from .bet_codec import encode_bet, decode_bet
from .opcode import Opcode

class BetPacket(Packet):
    def __init__(self, bets: list[Bet], agency_id: int):
        self.bets = bets
        self.agency_id = agency_id

    def to_bytes(self) -> bytes:
        payload = self.agency_id.to_bytes(1, byteorder='big')
        for bet in self.bets:
            payload += encode_bet(bet)
        return build_frame(Opcode.BET, payload)

def bet_pkt_from_bytes(data: bytes) -> Packet:
    agency_id = data[0]
    offset = 1
    bets = []

    while offset < len(data):
        bet, consumed = decode_bet(data[offset:], agency_id)
        bets.append(bet)
        offset += consumed

    return BetPacket(bets, agency_id)