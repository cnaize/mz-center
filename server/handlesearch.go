package server

import (
	"fmt"
	"github.com/cnaize/mz-common/log"
	"github.com/cnaize/mz-common/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) handleGetSearchRequestList(c *gin.Context) {
	db := s.config.DB

	var in struct {
		Offset uint `form:"offset"`
		Count  uint `form:"count"`
	}

	c.ShouldBindQuery(&in)
	if in.Count == 0 || in.Count >= model.MaxRequestItemsPerRequestCount {
		in.Count = model.MaxRequestItemsPerRequestCount
	}

	res, err := db.GetSearchRequestList(in.Offset, in.Count)
	if err != nil {
		if db.IsSearchItemNotFound(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, model.SearchRequestList{
			Error: &model.Error{Str: fmt.Sprintf("db failed: %+v", err)},
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) handleAddSearchRequest(c *gin.Context) {
	db := s.config.DB

	var inRequest model.SearchRequest
	if err := c.ShouldBindQuery(&inRequest); err != nil {
		log.Warn("Server: search request add failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := db.AddSearchRequest(inRequest); err != nil {
		log.Warn("Server: search request add failed: %+v", err)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	c.JSON(http.StatusCreated, inRequest)
}

func (s *Server) handleGetSearchResponseList(c *gin.Context) {
	db := s.config.DB

	var in struct {
		Offset uint `form:"offset"`
		Count  uint `form:"count"`
	}

	c.ShouldBindQuery(&in)
	if in.Count == 0 || in.Count >= model.MaxResponseItemsCount {
		in.Count = model.MaxResponseItemsCount
	}

	var inRequest model.SearchRequest
	if err := c.ShouldBindQuery(&inRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.SearchResponseList{
			Error: &model.Error{Str: fmt.Sprintf("input parse failed: %+v", err)},
		})
		return
	}

	res, err := db.GetSearchResponseList(inRequest, in.Offset, in.Count)
	if err != nil {
		if s.config.DB.IsSearchItemNotFound(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, model.SearchResponseList{
			Error: &model.Error{Str: fmt.Sprintf("db failed: %+v", err)},
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) handleAddSearchResponseList(c *gin.Context) {
	db := s.config.DB

	var inRequest model.SearchRequest
	if err := c.ShouldBindQuery(&inRequest); err != nil {
		log.Warn("Server: search response list add failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var inResponseList model.SearchResponseList
	if err := c.ShouldBindJSON(&inResponseList); err != nil {
		log.Warn("Server: search response list add failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	username := c.MustGet("username").(string)
	user, err := db.GetUser(model.User{Username: username})
	if err != nil {
		log.Error("Server: search response list add failed: %+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := db.AddSearchResponseList(user, inRequest, inResponseList); err != nil {
		if db.IsSearchItemNotFound(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Warn("Server: search response list add failed: %+v", err)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	c.JSON(http.StatusCreated, nil)
}
