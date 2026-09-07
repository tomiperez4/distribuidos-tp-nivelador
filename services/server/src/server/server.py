import socket
import logger
from protocol import Protocol, EndOfBets
from lottery import Lottery


class Server:
    def __init__(self, server_host: str, server_port: int, storage_path: str) -> None:
        self.server_host = server_host
        self.server_port = server_port
        self.lottery = Lottery(storage_path)

    def _handle_client(self, protocol: Protocol) -> None:
        action = "handle-client"
        while True:
            try:
                bet = protocol.recv_bet()
                self.lottery.store_bets([bet])


            except EndOfBets:
                break
            except Exception as e:
                logger.error(
                    action, logger.LogResult.fail
                )
                raise e

        bets = self.lottery.load_bets()

        for bet in bets:
            if self.lottery.has_won(bet):
                protocol.send_bet(bet)

        protocol.send_fin()

    def run(self):
        action = "accept-connection"
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_socket:
            server_socket.bind((self.server_host, self.server_port))
            server_socket.listen()
            while True:
                try:
                    logger.info(action, logger.LogResult.in_progress)
                    client_socket, _ = server_socket.accept()
                except Exception as e:
                    logger.error(action, logger.LogResult.fail)
                    raise e
                logger.info(action, logger.LogResult.success)
                protocol = Protocol(client_socket)

                self._handle_client(protocol)
