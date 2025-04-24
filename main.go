package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bdmoleg/go-go-project/app"
	"github.com/bdmoleg/go-go-project/configs"
	"github.com/bdmoleg/go-go-project/utils"
)

func main() {

	var extractInterval int
	flag.IntVar(&extractInterval, "interval", 30, "interval time for extraction in seconds")
	flag.Parse()

	opts := utils.PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := utils.NewPrettyHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	logger.Info("starting app")

	logger.Info("injecting env variables into process")
	configs.InitEnvVariables()

	logger.Info("app name", slog.String("AppName", configs.AppName))

	httpCLient := &http.Client{
		Timeout: configs.HttpTimeout,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	asanaExporter := app.NewAsanaExtractor(httpCLient, logger)

	// numOfGoroutines := 250
	// wg := &sync.WaitGroup{}
	// wg.Add(numOfGoroutines)
	// asanaExporter.ThreadsTestRateLimiting(wg, numOfGoroutines)
	// wg.Wait()

	asanaExporter.Extract(ctx, time.Duration(extractInterval)*time.Second)

}
