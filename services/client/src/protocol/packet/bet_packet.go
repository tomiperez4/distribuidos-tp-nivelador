package packet

import (
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/lottery"
)

type BetPacket struct {
	Bets []lottery.Bet
}

func NewBetPacket(bets []lottery.Bet) Packet {
	return &BetPacket{Bets: bets}
}

func betFromBytes(payload []byte) (Packet, error) {
	offset := 0
	var bets []lottery.Bet
	for offset < len(payload) {
		bet, consumed, err := deserializeBet(payload[offset:])
		if err != nil {
			return nil, err
		}
		offset += consumed
		bets = append(bets, bet)
	}

	return NewBetPacket(bets), nil
}

func (bPkt *BetPacket) ToBytes() ([]byte, error) {
	var buff []byte

	for _, bet := range bPkt.Bets {
		payload, err := serializeBet(bet)
		if err != nil {
			return nil, err
		}
		buff = append(buff, payload...)
	}

	return buildFrame(BET_OPCODE, buff), nil
}
