package client

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/lottery"
)

const _CSV_FIELD_COUNT = 5

func FromCsv(csvBetLine string) (lottery.Bet, error) {
	args := strings.Split(csvBetLine, ",")
	if len(args) != _CSV_FIELD_COUNT {
		return lottery.Bet{}, fmt.Errorf("decode csv bet: expected %d fields, got %d: %q", _CSV_FIELD_COUNT, len(args), csvBetLine)
	}

	firstName, lastName := args[0], args[1]
	document, err := strconv.ParseUint(args[2], 10, 32)
	if err != nil {
		return lottery.Bet{}, err
	}
	birthdate := args[3]

	number, err := strconv.ParseUint(args[4], 0, 32)
	if err != nil {
		return lottery.Bet{}, err
	}

	return lottery.Bet{
		FirstName: firstName,
		LastName:  lastName,
		Document:  uint32(document),
		Birthdate: birthdate,
		Number:    uint32(number),
	}, nil
}

func ToCsv(bet lottery.Bet) string {
	return fmt.Sprintf("%s,%s,%d,%s,%d\n",
		bet.FirstName, bet.LastName, bet.Document, bet.Birthdate, bet.Number)
}
