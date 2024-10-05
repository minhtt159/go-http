package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin/v2"
)

var (
	app = kingpin.New(filepath.Base(os.Args[0]), "Simple HTTP Server")

	verbose = app.Flag("verbose", "Verbose level").Short('v').Bool()
	server  = app.Flag("server", "Listen interface").Short('s').Default("127.0.0.1").IP()
	port    = app.Flag("port", "Listen port").Short('p').Default("8080").String()

	action = app.Command("serve", "Start in server mode").Default()
	path   = action.Arg("path", "Listen path").Default("/").String()
)

func main() {
	fmt.Println("All Args", os.Args)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case action.FullCommand():
		fmt.Println("Start in serve mode")
		fmt.Printf("All flags: \nverbose: %t\nserver: %s\nport: %s\npath: %s\n",
			*verbose, *server, *port, *path)
	default:
		log.Fatal("Action not valid")
	}
}
