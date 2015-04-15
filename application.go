package main

import (
	"application/server"
	"github.com/alecthomas/kingpin"
	"os"
)

var (
	app      = kingpin.New("application", "RESTful web API application")
	cmdPrint = app.Command("print-routes", "Print all supported routes")
	command  = kingpin.MustParse(app.Parse(os.Args[1:]))
)

func main() {
	switch command {
	case cmdPrint.FullCommand():
		server.PrintRoutes(server.Routes())
	default:
		app.Usage(os.Stdout)
		server.Start()
		server.Stop()
	}
}
