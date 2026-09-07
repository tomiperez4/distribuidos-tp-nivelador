package packet

import (
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/lottery"
)

type BetPacket struct {
	bet lottery.Bet
}

func NewBetPacket(bet lottery.Bet) Packet {
	return &BetPacket{
		bet: bet,
	}
}

func betFromBytes(payload []byte) (Packet, error) {
	parsedBet, err := decodeBet(payload)
	if err != nil {
		return nil, err
	}

	return NewBetPacket(parsedBet), nil
}

func (bPkt *BetPacket) ToBytes() ([]byte, error) {
	payload, err := encodeBet(bPkt.bet)
	if err != nil {
		return nil, err
	}

	return buildFrame(BET_OPCODE, payload), nil
}

func (bPkt *BetPacket) Bet() lottery.Bet {
	return bPkt.bet
}
