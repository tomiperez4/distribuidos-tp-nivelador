import socket
import signal
import logger
from protocol import Protocol
from lottery import Lottery
from multiprocessing import Process, Barrier
from .client_handler import ClientHandler

class Server:
    def __init__(self, server_host: str, server_port: int, storage_path: str, agency_min_quorum: int) -> None:
        self.server_host = server_host
        self.server_port = server_port
        self.lottery = Lottery(storage_path)
        self.barrier = Barrier(agency_min_quorum)
        self.workers = []
        self.server_socket = None
        self.shutting_down = False

    def _request_shutdown(self):
        self.shutting_down = True
        if self.server_socket:
            self.server_socket.close()

    def _shutdown(self):
        self.barrier.abort()

        for worker in self.workers:
            if worker.is_alive():
                worker.terminate()
                worker.join()

    def _handle_incoming(self):
        action = "accept-connection"
        while True:
            try:
                logger.info(action, logger.LogResult.in_progress)
                client_socket, _ = self.server_socket.accept()

            except OSError:
                if self.shutting_down:
                    return
                logger.error(action, logger.LogResult.fail)
                raise

            logger.info(action, logger.LogResult.success)

            handler = ClientHandler(Protocol(client_socket), self.lottery, self.barrier)
            worker = Process(target=handler.run)
            worker.start()
            client_socket.close()

            self.workers.append(worker)
            self.workers = [w for w in self.workers if w.is_alive()]


    def run(self):
        signal.signal(signal.SIGTERM, lambda signum, frame: self._request_shutdown())

        try:
            self.server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.server_socket.bind((self.server_host, self.server_port))
            self.server_socket.listen()

            self._handle_incoming()
        finally:
            self._shutdown()
