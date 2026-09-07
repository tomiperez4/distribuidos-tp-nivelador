package client

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/lottery"
)

func FromCsv(csvBetLine string, agencyId uint8) (lottery.Bet, error) {
	args := strings.Split(csvBetLine, ",")
	firstName, lastName := args[0], args[1]
	document, err := strconv.ParseUint(args[2], 10, 32)
	if err != nil {
		// Manejo error
		return lottery.Bet{}, err
	}
	birthdate := args[3]

	number, err := strconv.ParseUint(args[4], 0, 32)
	if err != nil {
		// Manejo error
		return lottery.Bet{}, err
	}

	return lottery.Bet{
		AgencyId:  agencyId,
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
