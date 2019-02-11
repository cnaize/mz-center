package server

import "github.com/cnaize/mz-center/db"

type Config struct {
	Port             uint
	MinMZCoreVersion string
	JwtTokenPassword string
	DB               db.DB
}
