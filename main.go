package main

import (
	"flag"
	"github.com/cnaize/mz-center/db/sqlite"
	"github.com/cnaize/mz-center/server"
	"github.com/cnaize/mz-common/log"
	"os"
	"strconv"
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

	if serverConfig.MinMZCoreVersion == "" {
		serverConfig.MinMZCoreVersion = os.Getenv("MIN_MZ_CORE_VERSION")
	}
	if serverConfig.JwtTokenPassword == "" {
		serverConfig.JwtTokenPassword = os.Getenv("JWT_TOKEN_PASSWORD")
	}
	if len(os.Getenv("PORT")) > 0 {
		port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 23)
		if err != nil {
			log.Fatal("MuzeZone Center: port parsing failed: %+v", err)
		}
		serverConfig.Port = uint(port)
	}

	db, err := sqlite.New()
	if err != nil {
		log.Fatal("MuzeZone Center: db open failed: %+v", err)
	}
	serverConfig.DB = db

	if err := server.New(serverConfig).Run(); err != nil {
		log.Fatal("MuzeZone Center: server run failed: %+v", err)
	}
}
