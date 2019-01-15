package handle

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListSearches(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"not": "implemented",
	})
}

func GetSearch(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"not": "implemented",
	})
}

func AddSearch(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"not": "implemented",
	})
}

func AddSearchResult(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"not": "implemented",
	})
}
