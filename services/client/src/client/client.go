package client

import (
	"bufio"
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/lottery"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol"
)

const CONNECTION_ATTEMPTS_MAX = 3
const CONNECTION_ATTEMPS_DELAY_MS = 200

type ClientConfig struct {
	ServerHost string
	ServerPort string
	AgencyId   uint8
	InputPath  string
	OutputPath string
	BatchSize  int
}

type Client struct {
	conn         protocol.Protocol
	config       ClientConfig
	shuttingDown atomic.Bool
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
	const mainAction = "client-running"

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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		client.shuttingDown.Store(true)
		client.conn.Close()
	}()

	if err := client.conn.Handshake(client.config.AgencyId); err != nil {
		return client.handleConnErr("handshake", err)
	}
	logger.Info("handshake", logger.Success, "agency-id", client.config.AgencyId)

	scanner := bufio.NewScanner(inFile)

	if action, err := client.communicationProcess(scanner, outFile); err != nil {
		return client.handleConnErr(action, err)
	}
	logger.Info(mainAction, logger.Success, "agency-id", client.config.AgencyId)

	return nil
}

func (client *Client) communicationProcess(scanner *bufio.Scanner, outFile *os.File) (string, error) {
	for {
		bets, err := readBatch(scanner, client.config.BatchSize)
		if err != nil {
			return "read-batch", err
		}

		if len(bets) == 0 {
			break
		}

		if err := client.conn.SendBets(bets); err != nil {
			return "send-message", err
		}

		if err := client.conn.RecvAck(); err != nil {
			return "recv-ack", err
		}
	}

	if err := client.conn.SendFin(); err != nil {
		return "send-fin", err
	}

	winners, err := client.conn.RecvWinners()
	if err != nil {
		return "recv-winners", err
	}

	for _, bet := range winners {
		if _, err := outFile.Write([]byte(ToCsv(bet))); err != nil {
			return "write-winner", err
		}
	}

	if err := outFile.Sync(); err != nil {
		return "sync-output-file", err
	}

	return "", nil
}

func (client *Client) handleConnErr(action string, err error) error {
	if client.shuttingDown.Load() && errors.Is(err, net.ErrClosed) {
		logger.Info(action, logger.InProgress, "reason", "shutting down")
		return nil
	}
	logger.Error(action, logger.Fail)
	return err
}

func readBatch(scanner *bufio.Scanner, size int) ([]lottery.Bet, error) {
	const mainAction = "read-batch"
	bets := make([]lottery.Bet, 0, size)

	for len(bets) < size && scanner.Scan() {
		bet, err := FromCsv(scanner.Text())
		if err != nil {
			return nil, err
		}
		bets = append(bets, bet)
	}

	return bets, nil
}
