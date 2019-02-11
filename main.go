package main

import (
	"flag"
	"github.com/cnaize/mz-center/db/sqlite"
	"github.com/cnaize/mz-center/server"
	"github.com/cnaize/mz-common/log"
)

const (
	MinMZCoreVersion = ""
	JwtTokenPassword = ""
)

var (
	loggerConfig log.Config
	serverConfig server.Config
)

func init() {
	serverConfig.MinMZCoreVersion = MinMZCoreVersion
	serverConfig.JwtTokenPassword = JwtTokenPassword

	flag.UintVar(&loggerConfig.Lvl, "log-lvl", 5, "log level")

	flag.UintVar(&serverConfig.Port, "port", 11310, "server port")
}

func main() {
	flag.Parse()
	log.Init(loggerConfig)

	db, err := sqlite.New()
	if err != nil {
		log.Fatal("MuzeZone Center: db open failed: %+v", err)
	}
	serverConfig.DB = db

	if err := server.New(serverConfig).Run(); err != nil {
		log.Fatal("MuzeZone Center: server run failed: %+v", err)
	}
}
