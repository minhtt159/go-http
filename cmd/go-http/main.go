package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin/v2"
)

var (
	verbose = kingpin.Flag("verbose", "Verbose level").Short('v').Bool()
	server  = kingpin.Arg("server", "Listen interface").String()
	port    = kingpin.Arg("port", "Listen port").String()
	path    = kingpin.Arg("path", "Listen path").String()
)

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]), "Simple HTTP Server").UsageWriter(os.Stdout)
	// app.HelpFlag.Short('h')

	// Parse Command
	parsedCmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	fmt.Println(parsedCmd)
}
