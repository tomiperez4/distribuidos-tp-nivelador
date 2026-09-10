package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	client "github.com/7574-sistemas-distribuidos/tp-nivelador/src/client"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
)

func loadConfig() (client.ClientConfig, error) {
	rawAgencyId := os.Getenv("AGENCY_ID")
	if rawAgencyId == "" {
		return client.ClientConfig{}, errors.New("AGENCY_ID environment variable is required")
	}
	agencyId, err := strconv.ParseUint(rawAgencyId, 10, 8)
	if err != nil {
		return client.ClientConfig{}, fmt.Errorf("AGENCY_ID must be a valid single-byte (0-255) number : %w", err)
	}

	serverHost := os.Getenv("SERVER_HOST")
	if serverHost == "" {
		return client.ClientConfig{}, errors.New("SERVER_HOST environment variable is required")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		return client.ClientConfig{}, errors.New("SERVER_PORT environment variable is required")
	}

	inputPath := os.Getenv("INPUT_FILE")
	if inputPath == "" {
		return client.ClientConfig{}, errors.New("INPUT_FILE environment variable is required")
	}

	outputPath := os.Getenv("OUTPUT_FILE")
	if outputPath == "" {
		return client.ClientConfig{}, errors.New("OUTPUT_FILE environment variable is required")
	}

	rawBatchSize := os.Getenv("BATCH_SIZE")
	if rawBatchSize == "" {
		return client.ClientConfig{}, errors.New("BATCH_SIZE environment variable is required")
	}
	batchSize, err := strconv.Atoi(rawBatchSize)
	if err != nil {
		return client.ClientConfig{}, fmt.Errorf("BATCH_SIZE must be a valid number: %w", err)
	}
	if batchSize <= 0 {
		return client.ClientConfig{}, errors.New("BATCH_SIZE must be greater than 0")
	}

	return client.ClientConfig{
		ServerHost: serverHost,
		ServerPort: serverPort,
		AgencyId:   uint8(agencyId),
		InputPath:  inputPath,
		OutputPath: outputPath,
		BatchSize:  batchSize,
	}, nil
}

func run() int {
	config, err := loadConfig()
	if err != nil {
		logger.Error("load-config", logger.Fail, "err", err)
		return 1
	}

	client, err := client.NewClient(config)
	if err != nil {
		logger.Error("client-new", logger.Fail, "err", err)
		return 1
	}

	if err := client.Run(); err != nil {
		logger.Error("client-run", logger.Fail, "err", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
