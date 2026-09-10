import os
import sys
from dataclasses import dataclass

import logger
import server


@dataclass(frozen=True)
class ServerConfig:
    server_host: str
    server_port: int
    storage_path: str
    agency_quorum_min: int


def load_config() -> ServerConfig:
    return ServerConfig(
        server_host=os.environ["SERVER_HOST"],
        server_port=int(os.environ["SERVER_PORT"]),
        storage_path=os.environ.get("STORAGE_PATH", "storage.csv"),
        agency_quorum_min=int(os.environ["AGENCY_QUORUM_MIN"]),
    )


def main():
    logger.init()
    config = load_config()
    s = server.Server(
        config.server_host,
        config.server_port,
        config.storage_path,
        config.agency_quorum_min,
    )
    try:
        s.run()
    except Exception as e:
        logger.error("server-run", logger.LogResult.fail, "err", e)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
