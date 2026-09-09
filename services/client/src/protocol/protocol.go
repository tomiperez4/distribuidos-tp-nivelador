package protocol

import (
	"errors"
	"net"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/lottery"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/packet"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

type Protocol struct {
	conn net.Conn
}

func NewProtocol(conn net.Conn) Protocol {

	return Protocol{conn: conn}
}

func (p *Protocol) SendBets(bets []lottery.Bet, agencyId uint8) error {
	betPacket := packet.NewBetPacket(bets, agencyId)
	return p.sendPkt(betPacket)
}

func (p *Protocol) RecvAck() error {
	pkt, err := p.recvPkt()
	if err != nil {
		return err
	}

	if _, ok := pkt.(*packet.AckPacket); !ok {
		return errors.New("expected ack packet")
	}

	return nil
}

func (p *Protocol) RecvWinners() ([]lottery.Bet, error) {
	winners := make([]lottery.Bet, 0)

	for {
		pkt, err := p.recvPkt()
		if err != nil {
			return nil, err
		}

		switch v := pkt.(type) {
		case *packet.BetPacket:
			winners = append(winners, v.Bets...)
		case *packet.FinPacket:
			return winners, nil
		default:
			return nil, errors.New("unkwnown packet type")
		}
	}
}

func (p *Protocol) SendFin() error {
	return p.sendPkt(packet.NewFinPacket())
}

func (p *Protocol) sendPkt(pkt packet.Packet) error {
	frame, err := pkt.ToBytes()
	if err != nil {
		return err
	}
	return safe_socket.SendAll(p.conn, frame)
}

func (p *Protocol) recvPkt() (packet.Packet, error) {
	header, err := safe_socket.RecvAll(p.conn, packet.HEADER_LENGTH)
	if err != nil {
		return nil, err
	}

	opCode, payloadLen, err := packet.ParseHeader(header)
	if err != nil {
		return nil, err
	}

	payload, err := safe_socket.RecvAll(p.conn, int(payloadLen))
	if err != nil {
		return nil, err
	}

	return packet.FromBytes(opCode, payload)
}

func (p *Protocol) Close() {
	p.conn.Close()
}
