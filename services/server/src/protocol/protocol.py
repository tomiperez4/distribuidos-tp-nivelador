import socket

from .packet import FinPacket, BetPacket, Packet, HEADER_LENGTH, parse_header, from_bytes
from lottery import Bet
from safe_socket import safe_socket

class EndOfBets(Exception):
    """all bets have been received"""

class Protocol:
    def __init__(self, client_socket: socket.socket):
        self.socket = client_socket

    def send_bet(self, bet: Bet) -> None:
        self.__send_pkt__(BetPacket([bet], bet.agency_id))

    def send_fin(self) -> None:
        self.__send_pkt__(FinPacket())

    def recv_batch(self) -> list[Bet]:
        packet = self.__recv_pkt__()

        if isinstance(packet, BetPacket):
            return packet.bets
        elif isinstance(packet, FinPacket):
            raise EndOfBets
        else:
            raise Exception("unknown packet type")

    def close(self) -> None:
        self.socket.close()


    def __send_pkt__(self, pkt: Packet) -> None:
        frame = pkt.to_bytes()
        return safe_socket.send_all(self.socket, frame)

    def __recv_pkt__(self) -> Packet:
        buff = safe_socket.recv_all(self.socket, HEADER_LENGTH)

        header = parse_header(buff)
        payload = safe_socket.recv_all(self.socket, header.length)

        return from_bytes(header.opcode, payload)