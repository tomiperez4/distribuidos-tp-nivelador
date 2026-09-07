import socket

def recv_all(socket: socket.socket, size):
    total = 0
    data_buffer = b''

    while total < size:
        data = socket.recv(size - total)
        if not data:
            if not data_buffer:
                return b""
            raise ConnectionError
        data_buffer += data
        total += len(data)

    return data_buffer


def send_all(socket: socket.socket, bytes):
    total = 0

    while total < len(bytes):
        amount = socket.send(bytes[total:])
        if not amount:
            raise ConnectionError
        total += amount

    return total
