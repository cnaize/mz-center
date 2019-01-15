package server

import (
	"fmt"
	"github.com/cnaize/mz-center/db"
	"github.com/cnaize/mz-center/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (s *Server) handleGetSearchRequestList(c *gin.Context) {
	type In struct {
		Offset int `form:"offset"`
		Count  int `form:"count"`
	}
	var in In

	if err := c.ShouldBindQuery(&in); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("parsing query string failed: %+v", err),
		})
		return
	}

	res, err := s.config.DB.GetSearchRequestList(in.Offset, in.Count)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("db failed: %+v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requestList": res,
	})
}

func (s *Server) handleAddSearchRequest(c *gin.Context) {
	text := ParseSearchText(c.Param("text"))
	if text == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid search text"),
		})
		return
	}

	if err := s.config.DB.AddSearchRequest(model.SearchRequest{Text: text, CreatedAt: time.Now()}); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("db failed: %+v", err),
		})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (s *Server) handleGetSearchResponseList(c *gin.Context) {
	text := ParseSearchText(c.Param("text"))
	if text == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid search text"),
		})
		return
	}

	type In struct {
		Offset int `form:"offset"`
		Count  int `form:"count"`
	}
	var in In

	if err := c.ShouldBindQuery(&in); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("parsing query string failed: %+v", err),
		})
		return
	}

	res, err := s.config.DB.GetSearchResponseList(model.SearchRequest{Text: text}, in.Offset, in.Count)
	if err != nil {
		if err == db.ErrorNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, nil)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("db failed: %+v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"responseList": res,
	})
}

func (s *Server) handleAddSearchResponseList(c *gin.Context) {
	text := ParseSearchText(c.Param("text"))
	if text == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid search text"),
		})
		return
	}

	username := c.Param("username")
	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid username"),
		})
		return
	}

	var resp model.SearchResponseList
	if err := c.BindJSON(&resp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid response list: %+v", err),
		})
		return
	}

	if err := s.config.DB.AddSearchResponseList(model.User{Username: username}, model.SearchRequest{Text: text}, resp); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("db failed: %+v", err),
		})
		return
	}

	c.JSON(http.StatusCreated, nil)
}
