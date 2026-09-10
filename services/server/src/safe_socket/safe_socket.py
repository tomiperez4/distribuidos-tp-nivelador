import socket

def recv_all(socket: socket.socket, size):
    chunks = []
    received = 0

    while received < size:
        data = socket.recv(size - received)

        if not data:
            return b""
        chunks.append(data)
        received += len(data)

    return b"".join(chunks)


def send_all(socket: socket.socket, bytes):
    sent_bytes = 0

    while sent_bytes < len(bytes):
        sent_bytes += socket.send(bytes[sent_bytes:])
