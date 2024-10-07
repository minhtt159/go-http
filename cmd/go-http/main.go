package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
)

func main() {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name(filepath.Base(os.Args[0])),
		kong.Description("Simple HTTP Server"),
	)
	logger := initLogger(cli)

	// Print context if DEBUG
	logger.Debug(fmt.Sprintf("All Args: %s", os.Args))
	logger.Debug(fmt.Sprintf("Context: %s", cli))

	// Start commands
	switch ctx.Command() {
	case "serve":
		logger.Info("Start in serve mode")
		logger.Debug(fmt.Sprintf("All flags: \nserver: %s\nport: %s\npath: %s\n",
			cli.Serve.Server, cli.Serve.Port, cli.Serve.Path))
	case "get <url>":
		logger.Info("Start in GET mode")
		logger.Debug(fmt.Sprintf("All flags:\nurl: %s\n", cli.Get.Url))
	default:
		logger.Error(fmt.Sprintf("Context: %s", cli))
		panic(ctx.Command())
	}
}
