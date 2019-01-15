package server

import (
	"github.com/cnaize/mz-center/log"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleLog(c *gin.Context) {
	log.Info("%s %s", c.Request.Method, c.Request.URL.String())
}
