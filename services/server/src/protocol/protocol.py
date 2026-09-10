import socket

from .packet import AckPacket, FinPacket, BetPacket, HelloPacket, Packet, HEADER_LENGTH, parse_header, from_bytes
from lottery import Bet
from safe_socket import safe_socket

class EndOfBets(Exception):
    """all bets have been received"""

class ConnectionClosed(Exception):
    """the peer closed the connection without sending a Fin packet"""

class Protocol:
    def __init__(self, client_socket: socket.socket):
        self.socket = client_socket
        self.agency_id: int | None = None

    def recv_hello(self) -> int:
        packet = self.__recv_pkt__()
        if not isinstance(packet, HelloPacket):
            raise Exception("expected hello packet")
        self.agency_id = packet.agency_id
        return packet.agency_id

    def send_fin(self) -> None:
        self.__send_pkt__(FinPacket())

    def send_ack(self) -> None:
        self.__send_pkt__(AckPacket())

    def send_winners(self, bets: list[Bet]) -> None:
        if not bets:
            return
        self.__send_pkt__(BetPacket(bets))

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
        if not buff:
            raise ConnectionClosed

        header = parse_header(buff)
        payload = safe_socket.recv_all(self.socket, header.length)
        if not payload and header.length > 0:
            raise ConnectionClosed

        return from_bytes(header.opcode, payload, self.agency_id)
