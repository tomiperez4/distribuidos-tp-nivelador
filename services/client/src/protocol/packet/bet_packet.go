package packet

import (
	"errors"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/lottery"
)

const _BET_PACKET_HEADER_SIZE = 1

type BetPacket struct {
	AgencyId uint8
	Bets     []lottery.Bet
}

func NewBetPacket(bets []lottery.Bet, agencyId uint8) Packet {
	return &BetPacket{
		AgencyId: agencyId,
		Bets:     bets,
	}
}

func betFromBytes(payload []byte) (Packet, error) {
	if len(payload) < _BET_PACKET_HEADER_SIZE {
		return nil, errors.New("decode bet packet: payload too short for header")
	}

	agencyId := payload[0]
	offset := _BET_PACKET_HEADER_SIZE
	var bets []lottery.Bet
	for offset < len(payload) {
		bet, consumed, err := deserializeBet(payload[offset:])
		if err != nil {
			return nil, err
		}
		offset += consumed
		bets = append(bets, bet)
	}

	return NewBetPacket(bets, agencyId), nil
}

func (bPkt *BetPacket) ToBytes() ([]byte, error) {
	buff := []byte{bPkt.AgencyId}

	for _, bet := range bPkt.Bets {
		payload, err := serializeBet(bet)
		if err != nil {
			return nil, err
		}
		buff = append(buff, payload...)
	}

	return buildFrame(BET_OPCODE, buff), nil
}
