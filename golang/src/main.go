package main

import (
	"os"
	"runtime"
	"server"
)

var version string

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//serverType := os.Args[1]
	config := os.Args[2]

	server.Main(config, version)
}
