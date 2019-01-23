package server

import (
	"fmt"
	"github.com/cnaize/mz-center/db"
	"github.com/cnaize/mz-common/log"
	"github.com/cnaize/mz-common/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) handleGetSearchRequestList(c *gin.Context) {
	var in struct {
		Offset uint `form:"offset"`
		Count  uint `form:"count"`
	}

	if err := c.ShouldBindQuery(&in); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.SearchRequestList{
			Error: &model.Error{Str: fmt.Sprintf("query string parse failed: %+v", err)},
		})
		return
	}

	res, err := s.config.DB.GetSearchRequestList(in.Offset, in.Count)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.SearchRequestList{
			Error: &model.Error{Str: fmt.Sprintf("db failed: %+v", err)},
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) handleAddSearchRequest(c *gin.Context) {
	text := ParseSearchText(c.Param("text"))
	if text == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.SearchRequest{
			Error: &model.Error{Str: fmt.Sprintf("invalid search text")},
		})
		return
	}

	log.Debug("Search text: %s", text)

	req, err := s.config.DB.AddSearchRequest(&model.SearchRequest{Text: text})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.SearchRequest{
			Error: &model.Error{Str: fmt.Sprintf("db failed: %+v", err)},
		})
		return
	}

	c.JSON(http.StatusCreated, req)
}

func (s *Server) handleGetSearchResponseList(c *gin.Context) {
	var in struct {
		Offset uint `form:"offset"`
		Count  uint `form:"count"`
	}

	text := ParseSearchText(c.Param("text"))
	if text == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.SearchResponseList{
			Error: &model.Error{Str: fmt.Sprintf("invalid search text")},
		})
		return
	}

	if err := c.ShouldBindQuery(&in); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.SearchResponseList{
			Error: &model.Error{Str: fmt.Sprintf("query string parse failed: %+v", err)},
		})
		return
	}

	res, err := s.config.DB.GetSearchResponseList(&model.SearchRequest{Text: text}, in.Offset, in.Count)
	if err != nil {
		if err != db.ErrorNotFound {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.SearchResponseList{
				Error: &model.Error{Str: fmt.Sprintf("db failed: %+v", err)},
			})
			return
		}

		res = &model.SearchResponseList{Request: &model.SearchRequest{}}
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) handleAddSearchResponseList(c *gin.Context) {
	text := ParseSearchText(c.Param("text"))
	if text == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.SearchResponseList{
			Error: &model.Error{Str: fmt.Sprintf("invalid search text")},
		})
		return
	}

	username := ParseSearchText(c.Param("username"))
	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.SearchResponseList{
			Error: &model.Error{Str: fmt.Sprintf("invalid username")},
		})
		return
	}

	resp := &model.SearchResponseList{}
	if err := c.BindJSON(resp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.SearchResponseList{
			Error: &model.Error{Str: fmt.Sprintf("invalid response list: %+v", err)},
		})
		return
	}

	if err := s.config.DB.AddSearchResponseList(&model.User{Username: username}, &model.SearchRequest{Text: text}, resp); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.SearchResponseList{
			Error: &model.Error{Str: fmt.Sprintf("db failed: %+v", err)},
		})
		return
	}

	c.JSON(http.StatusCreated, nil)
}
