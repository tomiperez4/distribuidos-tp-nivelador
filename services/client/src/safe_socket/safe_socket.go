package safe_socket

import "io"

//TODO: Complete with a short-read/short-write tolerant implementation

func SendAll(socket io.Writer, bytes []byte) error {
	total := 0

	for total < len(bytes) {
		amount, err := socket.Write(bytes[total:])
		if err != nil {
			return err
		}
		total += amount
	}
	return nil
}

func RecvAll(socket io.Reader, size int) ([]byte, error) {
	total := 0
	buff := make([]byte, size)

	for total < size {
		amount, err := socket.Read(buff[total:])
		if err != nil {
			return nil, err
		}
		total += amount
	}
	return buff, nil
}
