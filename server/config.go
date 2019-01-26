package server

import "github.com/cnaize/mz-center/db"

type Config struct {
	Port             uint
	MinMZCoreVersion string
	DB               db.DB
}
