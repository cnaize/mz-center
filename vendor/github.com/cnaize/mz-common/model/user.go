package model

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type User struct {
	Base
	Username  string     `json:"username" gorm:"unique_index"`
	PassHash  string     `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Error     *Error     `json:"error,omitempty" sql:"-" form:"error"`
}

type Token struct {
	jwt.StandardClaims
	Username string
}
