package server

import (
	"github.com/cnaize/mz-common/log"
	"github.com/cnaize/mz-common/model"
	"github.com/cnaize/mz-common/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) handleCreateUser(c *gin.Context) {
	var inUser model.User
	if err := c.ShouldBindJSON(&inUser); err != nil {
		log.Warn("Server: user creation failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if _, err := s.config.DB.GetUser(inUser); err == nil {
		log.Warn("Server: user creation failed: user already exists: %s", inUser.Username)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"xtoken": util.GenerateXToken(),
	})
}

func (s *Server) handleAuthCheck(c *gin.Context) {
	xtoken := c.GetHeader("X-Token")

	user, err := s.config.DB.GetUserByXToken(xtoken)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", &user)
}
