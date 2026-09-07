from .packet import Packet, build_frame
from lottery import Bet
from .bet_codec import encode_bet, decode_bet
from .opcode import Opcode

class BetPacket(Packet):
    def __init__(self, bet: Bet):
        self.bet = bet

    def to_bytes(self) -> bytes:
        return build_frame(Opcode.BET, encode_bet(self.bet))

def bet_pkt_from_bytes(data: bytes) -> Packet:
    return BetPacket(decode_bet(data))