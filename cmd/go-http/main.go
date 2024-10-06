package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/minhtt159/go-http/internal/logger"
)

// NOTE: Learn more about Kong: https://github.com/alecthomas/kong/blob/master/_examples/shell/commandstring/main.go
var CLI struct {
	Serve struct {
		Server string `arg:"" optional:"" help:"Listen Interface." default:"0.0.0.0"`
		Port   string `arg:"" optional:"" help:"Listen Port." default:"8080"`
		Path   string `arg:"" optional:"" help:"Listen Path." default:"/"`
	} `cmd:"Start in server mode" default:"withargs"`

	Get struct {
		url string `arg:"" help:"URL to GET"`
	} `cmd:"Get something"`

	Logging struct {
		Level string `enum:"debug,info,warn,error" default:"info"`
		Type  string `enum:"file,console" default:"console"`
	} `embed:"" prefix:"logging."`
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name(filepath.Base(os.Args[0])),
		kong.Description("Simple HTTP Server"),
	)

	// Initialize log handler
	var logLevel slog.Level
	err := logLevel.UnmarshalText([]byte(CLI.Logging.Level))
	if err != nil {
		fmt.Println("Unable to parge log level, default to Info")
		logLevel = slog.LevelInfo
	}

	logHandler := logger.New(
		&slog.HandlerOptions{Level: logLevel},
		logger.WithDestinationWriter(os.Stdout), logger.WithColor(),
	)
	logger := slog.New(logHandler)

	// Print context if DEBUG
	logger.Debug(fmt.Sprintf("All Args: %s", os.Args))
	logger.Debug(fmt.Sprintf("Context: %s", CLI))

	// Start commands
	switch ctx.Command() {
	case "serve":
		logger.Info("Start in serve mode")
		logger.Debug(fmt.Sprintf("All flags: \nserver: %s\nport: %s\npath: %s\n",
			CLI.Serve.Server, CLI.Serve.Port, CLI.Serve.Path))
	case "get":
		logger.Info("GET request")
	default:
		logger.Error(fmt.Sprintf("Context: %s", CLI))
		panic(ctx.Command())
	}
}
