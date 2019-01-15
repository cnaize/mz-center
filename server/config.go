package server

import "github.com/cnaize/mz-center/db"

type Config struct {
	Port uint
	DB   db.DB
}
