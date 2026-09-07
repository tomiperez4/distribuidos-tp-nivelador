package client

import (
	"bufio"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol"
)

const CONNECTION_ATTEMPTS_MAX = 3
const CONNECTION_ATTEMPS_DELAY_MS = 200

type ClientConfig struct {
	ServerHost string
	ServerPort string
	AgencyId   string
	InputPath  string
	OutputPath string
}

type Client struct {
	conn   protocol.Protocol
	config ClientConfig
}

func NewClient(config ClientConfig) (*Client, error) {
	conn, err := connectToServer(config.ServerHost, config.ServerPort)
	if err != nil {
		logger.Warn("connect-to-server", logger.Fail)
		return nil, err
	}

	betProtocol := protocol.NewProtocol(conn)

	client := &Client{conn: betProtocol, config: config}
	return client, nil
}

func connectToServer(host, port string) (net.Conn, error) {
	const action = "connect-to-server"
	var err error
	var conn net.Conn

	logger.Info(action, logger.InProgress)
	for i := range CONNECTION_ATTEMPTS_MAX {
		conn, err = net.Dial("tcp", host+":"+port)
		if err != nil {
			logger.Warn(action, logger.Fail, "attempt", i)
			time.Sleep(CONNECTION_ATTEMPS_DELAY_MS * time.Millisecond)
			continue
		}

		logger.Info(action, logger.Success)
		break
	}

	return conn, err
}

func (client *Client) Run() error {
	const mainAction = "test-echo-server"
	inFile, inErr := os.Open(client.config.InputPath)
	if inErr != nil {
		logger.Error("open-file", logger.Fail)
		return inErr
	}
	defer inFile.Close()

	outFile, outErr := os.Create(client.config.OutputPath)
	if outErr != nil {
		logger.Error("open-file", logger.Fail)
		return outErr
	}
	defer outFile.Close()

	defer client.conn.Close()

	scanner := bufio.NewScanner(inFile)
	var id uint8 = 0
	agencyId, _ := strconv.ParseUint(client.config.AgencyId, 10, 8)

	for scanner.Scan() {
		clientMessage := scanner.Text()
		messageArgs := []any{"agency-id", client.config.AgencyId, "message-id", id}

		bet, _ := FromCsv(clientMessage, uint8(agencyId))

		if err := client.conn.SendBet(bet); err != nil {
			logger.Error("send-message", logger.Fail, messageArgs...)
			return err
		}
	}

	if err := client.conn.SendFin(); err != nil {
		logger.Error("send-fin", logger.Fail)
		return err
	}

	for {
		bet, err := client.conn.RecvBet()
		if err != nil {
			if err == protocol.ErrEndOfWinners {
				break
			}
			logger.Error("recv-bet", logger.Fail)
			return err
		}
		outFile.Write([]byte(ToCsv(bet)))
	}
	logger.Info(mainAction, logger.Success, "agency-id", client.config.AgencyId)

	return nil
}
