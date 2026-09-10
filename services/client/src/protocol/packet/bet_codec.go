package packet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/lottery"
)

const (
	_BET_FIXED_FIELDS_SIZE = 14
	_DATE_BYTE_SIZE        = 4
)

func deserializeBet(data []byte) (lottery.Bet, int, error) {
	if len(data) < _BET_FIXED_FIELDS_SIZE {
		return lottery.Bet{}, 0, errors.New("decode bet: header too short")
	}

	lenFirst, lenLast := int(data[12]), int(data[13])
	total := _BET_FIXED_FIELDS_SIZE + lenFirst + lenLast
	if len(data) < total {
		return lottery.Bet{}, 0, errors.New("decode bet: short payload")
	}

	payloadLen := _BET_FIXED_FIELDS_SIZE
	firstName := string(data[payloadLen : payloadLen+lenFirst])
	payloadLen += lenFirst
	lastName := string(data[payloadLen : payloadLen+lenLast])

	return lottery.Bet{
		Document:  binary.BigEndian.Uint32(data[0:4]),
		Birthdate: unpackDate(data[4:8]),
		Number:    binary.BigEndian.Uint32(data[8:12]),
		FirstName: firstName,
		LastName:  lastName,
	}, total, nil
}

func serializeBet(bet lottery.Bet) ([]byte, error) {
	firstName, lastName := []byte(bet.FirstName), []byte(bet.LastName)
	if len(firstName) > math.MaxUint8 || len(lastName) > math.MaxUint8 {
		return nil, errors.New("encode bet: first or last name too long")
	}

	buff := make([]byte, 0, _BET_FIXED_FIELDS_SIZE+len(firstName)+len(lastName))

	buff = binary.BigEndian.AppendUint32(buff, bet.Document)
	buff = append(buff, packDate(bet.Birthdate)...)
	buff = binary.BigEndian.AppendUint32(buff, bet.Number)
	buff = append(buff, uint8(len(firstName)))
	buff = append(buff, uint8(len(lastName)))
	buff = append(buff, firstName...)
	buff = append(buff, lastName...)

	return buff, nil
}

func packDate(date string) []byte {
	year, month, day := parseDate(date)
	birthdate := year*10000 + month*100 + day

	buff := make([]byte, 0, _DATE_BYTE_SIZE)

	return binary.BigEndian.AppendUint32(buff, birthdate)
}

func unpackDate(data []byte) string {
	intDate := binary.BigEndian.Uint32(data)
	year, month, day := intDate/10000, (intDate/100)%100, intDate%100
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func parseDate(date string) (uint32, uint32, uint32) {
	dateArr := strings.Split(date, "-")
	year, errY := strconv.ParseUint(dateArr[0], 10, 16)
	month, errM := strconv.ParseUint(dateArr[1], 10, 8)
	day, errD := strconv.ParseUint(dateArr[2], 10, 8)

	if errY != nil || errM != nil || errD != nil {
		return 0, 0, 0
	}

	return uint32(year), uint32(month), uint32(day)
}
