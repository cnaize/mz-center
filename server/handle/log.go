package handle

import (
	"github.com/cnaize/mz-center/log"
	"github.com/gin-gonic/gin"
)

func Log(c *gin.Context) {
	log.Info("%s %s", c.Request.Method, c.Request.URL.String())
}
