package main

import (
	"flag"
	"github.com/cnaize/mz-center/db/memory"
	"github.com/cnaize/mz-center/server"
	"github.com/cnaize/mz-common/log"
)

var (
	loggerConfig log.Config
	serverConfig server.Config
)

func init() {
	flag.UintVar(&loggerConfig.Lvl, "log-lvl", 5, "log level")

	flag.UintVar(&serverConfig.Port, "port", 11310, "server port")
}

func main() {
	flag.Parse()
	log.Init(loggerConfig)

	serverConfig.DB = memory.NewDB()

	if err := server.New(serverConfig).Run(); err != nil {
		log.Fatal("MuzeZone Center: server run failed: %+v", err)
	}
}
