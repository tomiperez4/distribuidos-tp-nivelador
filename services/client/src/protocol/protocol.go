package protocol

import (
	"errors"
	"net"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/lottery"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/packet"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

type Protocol struct {
	conn net.Conn
}

var ErrEndOfWinners = errors.New("all winners have been received")

func NewProtocol(conn net.Conn) Protocol {

	return Protocol{conn: conn}
}

func (p *Protocol) SendBet(bet lottery.Bet) error {
	betPacket := packet.NewBetPacket(bet)
	return p.sendPkt(betPacket)
}

func (p *Protocol) RecvBet() (lottery.Bet, error) {
	pkt, err := p.recvPkt()
	if err != nil {
		return lottery.Bet{}, err
	}

	switch v := pkt.(type) {
	case *packet.BetPacket:
		return v.Bet(), nil
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
