package main

import (
	"client"
	"os"
	"runtime"
	"server"
)

var version string

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	runType := os.Args[1]
	confPath := os.Args[2]

	switch runType {
	case "test":
		client.PerfomanceTest(confPath)
	case "run":
		server.Main(confPath, "")
	}

}
