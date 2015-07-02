package main

import (
	"application/server"
	"flag"
)

var (
	boolPtr = flag.Bool("print-routes", false, "Print all supported routes")
)

func main() {
	flag.Parse()

	if *boolPtr {
		server.PrintRoutes(server.Routes())
	} else {
		server.Start()
		server.Stop()
	}
}
