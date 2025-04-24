package configs

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	AppName          string
	AsanaAPIEndpoint string
	AsanaToken       string
	HttpTimeout      time.Duration
)

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		os.Exit(1)
	}
}

func InitEnvVariables() {
	logger := slog.Default()

	AppName = os.Getenv("APP_NAME")
	if AppName == "" {
		AppName = "Asana"
	}

	AsanaAPIEndpoint = os.Getenv("ASANA_API_ENDPOINT")
	if AsanaAPIEndpoint == "" {
		AsanaAPIEndpoint = "https://app.asana.com/api/1.0"
	}

	AsanaToken = os.Getenv("ASANA_TOKEN")
	if AsanaToken == "" {
		logger.Error("asana token is empty")
		os.Exit(1)
	}

	timeoutValue := os.Getenv("HTTP_CLIENT_TIMEOUT")
	timeoutInt, err := strconv.Atoi(timeoutValue)
	if err != nil {
		logger.Error("timeout value is not a number",
			slog.String("error", err.Error()), slog.String("actualValue", timeoutValue))

		os.Exit(1)
	}
	HttpTimeout = time.Duration(timeoutInt) * time.Second
}
