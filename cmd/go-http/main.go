package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin/v2"
)

var (
	app = kingpin.New(filepath.Base(os.Args[0]), "Simple HTTP Server")

	verbose = app.Flag("verbose", "Verbose level").Short('v').Bool()
	server  = app.Flag("server", "Listen interface").Short('s').Default("0.0.0.0").IP()
	port    = app.Flag("port", "Listen port").Short('p').Default("8080").String()

	// Cmd: go-http serve
	action = app.Command("serve", "Start in server mode").Default()
	path   = action.Arg("path", "Listen path").Default("/").String()
)

func main() {
	log.Printf("All Args: %s", os.Args)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case action.FullCommand():
		log.Println("Start in serve mode")
		log.Printf("All flags: \nverbose: %t\nserver: %s\nport: %s\npath: %s\n",
			*verbose, *server, *port, *path)
	default:
		log.Fatal("Action not valid")
	}
}
