import signal
from collections.abc import Iterator
from multiprocessing import Barrier
from threading import BrokenBarrierError
import fcntl

import logger
from protocol import Protocol, EndOfBets, ConnectionClosed
from lottery import Lottery, Bet
from .file_lock import _file_lock


class ClientHandler:
    def __init__(self, protocol: Protocol, lottery: Lottery, barrier: Barrier) -> None:
        self.protocol = protocol
        self.lottery = lottery
        self.barrier = barrier
        self.shutting_down = False

    def run(self) -> None:
        signal.signal(signal.SIGTERM, self._request_shutdown)

        agency_id, completed = self._receive_bets()
        if not completed:
            return

        try:
            self.barrier.wait()
        except BrokenBarrierError:
            self.protocol.close()
            return

        winners = [
            bet for bet in self._load_bets_locked()
            if self.lottery.has_won(bet) and bet.agency_id == agency_id
        ]

        self.protocol.send_winners(winners)
        self.protocol.send_fin()
        self.protocol.close()

    def _request_shutdown(self) -> None:
        self.shutting_down = True
        self.protocol.close()

    def _receive_bets(self) -> tuple[int | None, bool]:
        action = "handle-client"
        agency_id = None

        while True:
            try:
                bets = self.protocol.recv_batch()
                if bets:
                    agency_id = bets[0].agency_id

                self._store_bets_locked(bets)

                self.protocol.send_ack()

            except EndOfBets:
                return agency_id, True
            except ConnectionClosed:
                logger.info(action, logger.LogResult.fail)
                self.protocol.close()
                return agency_id, False
            except OSError:
                if self.shutting_down:
                    return agency_id, False
                logger.error(action, logger.LogResult.fail)
                raise
            except Exception as e:
                logger.error(action, logger.LogResult.fail)
                raise e

    def _store_bets_locked(self, bets: list[Bet]) -> None:
        with _file_lock(self.lottery.storage_path, fcntl.LOCK_EX):
            self.lottery.store_bets(bets)

    def _load_bets_locked(self) -> Iterator[Bet]:
        with _file_lock(self.lottery.storage_path, fcntl.LOCK_SH):
            yield from self.lottery.load_bets()
