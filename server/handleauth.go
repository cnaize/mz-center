package server

import (
	"github.com/cnaize/mz-common/log"
	"github.com/cnaize/mz-common/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (s *Server) handleAuthCheck(c *gin.Context) {
	tHeader := c.GetHeader("Authorization")
	if tHeader == "" {
		log.Warn("Server: user with empty jwt token detected")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tBody := strings.Split(tHeader, " ")
	if len(tBody) != 2 {
		log.Warn("Server: bad user jwt token: %s", tHeader)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	var tk model.Token
	if t, err := jwt.ParseWithClaims(tBody[1], &tk, func(*jwt.Token) (interface{}, error) {
		return []byte(s.config.JwtTokenPassword), nil
	}); err != nil || !t.Valid {
		log.Warn("Server: invalid user jwt token: %+v", err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	log.Debug("Server: user %s: auth complete", tk.Username)

	c.Set("user", tk.Username)
}
