package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
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

	// log.Printf("All Args: %s", os.Args)
	// log.Printf("Context: %s", CLI)

	switch ctx.Command() {
	case "serve":
		log.Println("Start in serve mode")
		log.Printf("All flags: \nserver: %s\nport: %s\npath: %s\n",
			CLI.Serve.Server, CLI.Serve.Port, CLI.Serve.Path)
	case "get":
		log.Println("GET request")
	default:
		panic(ctx.Command())
	}
}
