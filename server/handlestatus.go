package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) handleStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
