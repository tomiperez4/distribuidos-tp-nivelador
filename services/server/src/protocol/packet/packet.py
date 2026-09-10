from abc import ABC, abstractmethod

HEADER_LENGTH = 5
AGENCY_ID_SIZE = 1

class Header:
    def __init__(self, opcode: int, length: int):
        self.opcode = opcode
        self.length = length

class Packet(ABC):
    @abstractmethod
    def to_bytes(self) -> bytes:
        pass

def build_frame(opcode: int, payload: bytes) -> bytes:
    return b"".join([
        opcode.to_bytes(1, byteorder="big"),
        len(payload).to_bytes(4, byteorder="big"),
        payload,
    ])

def parse_header(buff: bytes) -> Header:
    if len(buff) < HEADER_LENGTH:
        raise Exception("header too short")
    opcode = int.from_bytes(buff[:1], byteorder="big")
    length = int.from_bytes(buff[1:5], byteorder="big")

    return Header(opcode, length)