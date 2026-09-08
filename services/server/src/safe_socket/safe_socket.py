import socket

def recv_all(socket: socket.socket, size):
    data_buffer = b""

    while len(data_buffer) < size:
        data = socket.recv(size - len(data_buffer))

        if not data:
            return b""
        data_buffer += data

    return data_buffer


def send_all(socket: socket.socket, bytes):
    sent_bytes = 0

    while sent_bytes < len(bytes):
        sent_bytes += socket.send(bytes[sent_bytes:])
