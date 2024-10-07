package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/minhtt159/go-http/internal/logger"
)

// NOTE: Learn more about Kong: https://github.com/alecthomas/kong/blob/master/_examples/shell/commandstring/main.go
type CLI struct {
	Serve struct {
		Server string `arg:"" optional:"" short:"s" help:"Listen Interface." default:"0.0.0.0"`
		Port   string `arg:"" optional:"" short:"p" help:"Listen Port." default:"8080"`
		Path   string `arg:"" optional:"" short:"t" help:"Listen Path." default:"/"`
	} `cmd:"Start in server mode" default:"withargs"`

	Get struct {
		Url string `arg:"" short:"u" help:"URL to GET"`
	} `cmd:"Get something"`

	Logging struct {
		Level string `enum:"debug,info,warn,error" default:"info"`
		Type  string `enum:"file,console" default:"console"`
	} `embed:"" prefix:"logging."`
}

func initLogger(cli CLI) *slog.Logger {
	// Initialize log handler
	var logLevel slog.Level
	err := logLevel.UnmarshalText([]byte(cli.Logging.Level))
	if err != nil {
		fmt.Println("WARNING: Unable to parge log level, default to Info")
		logLevel = slog.LevelInfo
	}

	// Initialize log writer
	var logWriter io.Writer
	if cli.Logging.Type != "console" {
		fmt.Println("WARNING: Only support console log at the moment")
		logWriter = os.Stdout
	}

	logHandler := logger.New(
		&slog.HandlerOptions{Level: logLevel},
		logger.WithDestinationWriter(logWriter), logger.WithColor(),
	)
	logger := slog.New(logHandler)

	return logger
}
