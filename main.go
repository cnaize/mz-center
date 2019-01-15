package main

import (
	"flag"
	"github.com/cnaize/mz-center/log"
	"github.com/cnaize/mz-center/server"
)

var (
	loggerConfig log.Config
	serverConfig server.Config
)

func init() {
	flag.UintVar(&loggerConfig.Lvl, "log-lvl", 5, "log level")
	flag.StringVar(&loggerConfig.Dir, "log-dir", ".", "log directory")

	flag.UintVar(&serverConfig.Port, "port", 11312, "server port")
}

func main() {
	flag.Parse()
	log.Init(loggerConfig)

	if err := server.New(serverConfig).Run(); err != nil {
		log.Fatal("Server run failed: %+v", err)
	}
}
