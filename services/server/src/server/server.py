import socket
from collections.abc import Iterator

import logger
from protocol import Protocol, EndOfBets
from lottery import Lottery, Bet
from multiprocessing import Process, Barrier
import fcntl
from .file_lock import _file_lock


class Server:
    def __init__(self, server_host: str, server_port: int, storage_path: str, agency_min_quorum: int) -> None:
        self.server_host = server_host
        self.server_port = server_port
        self.lottery = Lottery(storage_path)
        self.barrier = Barrier(agency_min_quorum)
        self.workers = []

    @classmethod
    def _store_bets_locked(cls, bets: list[Bet], lottery: Lottery) -> None:
        with _file_lock(lottery.storage_path, fcntl.LOCK_EX):
            lottery.store_bets(bets)

    @classmethod
    def _load_bets_locked(cls, lottery: Lottery) -> Iterator[Bet]:
        with _file_lock(lottery.storage_path, fcntl.LOCK_SH):
            yield from lottery.load_bets()

    @classmethod
    def _handle_client(cls, protocol: Protocol, lottery: Lottery, barrier: Barrier) -> None:
        agency_id = None
        action = "handle-client"
        while True:
            try:
                bets = protocol.recv_batch()
                if bets:
                    agency_id = bets[0].agency_id

                cls._store_bets_locked(bets, lottery)

                protocol.send_ack()

            except EndOfBets:
                break
            except Exception as e:
                logger.error(
                    action, logger.LogResult.fail
                )
                raise e

        barrier.wait()

        winners = [
            bet for bet in cls._load_bets_locked(lottery)
            if lottery.has_won(bet) and bet.agency_id == agency_id
        ]

        protocol.send_winners(winners)

        protocol.send_fin()
        protocol.close()

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

                p = Process(target=self._handle_client, args=(protocol, self.lottery, self.barrier))
                p.start()
                client_socket.close()

                self.workers.append(p)
                self.workers = [w for w in self.workers if w.is_alive()]
