package protocol

import (
	"errors"
	"fmt"
	"net"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/lottery"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol/packet"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

type Protocol struct {
	conn net.Conn
}

var ErrEndOfWinners = errors.New("all winners have been received")

func NewProtocol(conn net.Conn) Protocol {

	return Protocol{conn: conn}
}

func (p *Protocol) SendBets(bets []lottery.Bet, agencyId uint8) error {
	betPacket := packet.NewBetPacket(bets, agencyId)
	return p.sendPkt(betPacket)
}

func (p *Protocol) RecvBet() (lottery.Bet, error) {
	pkt, err := p.recvPkt()
	if err != nil {
		return lottery.Bet{}, err
	}

	switch v := pkt.(type) {
	case *packet.BetPacket:
		if len(v.Bets) != 1 {
			return lottery.Bet{}, fmt.Errorf("Expected 1 bet only, but got %d", len(v.Bets))
		}
		return v.Bets[0], nil
	case *packet.FinPacket:
		return lottery.Bet{}, ErrEndOfWinners
	default:
		return lottery.Bet{}, errors.New("unkwnown packet type")
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
