package server

import (
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

		log.Error("Server: search request list get failed: %+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) handleAddSearchRequest(c *gin.Context) {
	db := s.config.DB
	username := c.MustGet("user").(string)

	var inRequest model.SearchRequest
	if err := c.ShouldBindQuery(&inRequest); err != nil {
		log.Warn("Server: search request add failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// TODO:
	//  handle "protected" mode

	if inRequest.Mode != model.MediaAccessTypePublic && inRequest.Mode != model.MediaAccessTypePrivate {
		inRequest.Mode = model.MediaAccessTypePublic
	}

	if err := db.AddSearchRequest(model.User{Username: username}, inRequest); err != nil {
		log.Warn("Server: search request add failed: %+v", err)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	c.JSON(http.StatusCreated, inRequest)
}

// NOTE:
//  returns 404 if related request not found
//  returns empty list if responses doesn't exists
func (s *Server) handleGetSearchResponseList(c *gin.Context) {
	db := s.config.DB
	username := c.MustGet("user").(string)

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
		log.Debug("Server: search response list get failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// TODO:
	//  handle "protected" mode

	if inRequest.Mode != model.MediaAccessTypePublic && inRequest.Mode != model.MediaAccessTypePrivate {
		inRequest.Mode = model.MediaAccessTypePublic
	}

	res, err := db.GetSearchResponseList(model.User{Username: username}, inRequest, in.Offset, in.Count)
	if err != nil {
		if db.IsSearchItemNotFound(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Error("Server: search response list get failed: +v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) handleAddSearchResponseList(c *gin.Context) {
	db := s.config.DB
	username := c.MustGet("user").(string)

	var inRequest model.SearchRequest
	if err := c.ShouldBindQuery(&inRequest); err != nil {
		log.Warn("Server: search response list add failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// TODO:
	//  handle "protected" mode

	if inRequest.Mode != model.MediaAccessTypePublic && inRequest.Mode != model.MediaAccessTypePrivate {
		inRequest.Mode = model.MediaAccessTypePublic
	}

	var inResponseList model.SearchResponseList
	if err := c.ShouldBindJSON(&inResponseList); err != nil {
		log.Warn("Server: search response list add failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := db.AddSearchResponseList(model.User{Username: username}, inRequest, inResponseList); err != nil {
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
