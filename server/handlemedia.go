package server

import (
	"github.com/cnaize/mz-common/log"
	"github.com/cnaize/mz-common/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) handleGetMediaRequestList(c *gin.Context) {
	db := s.config.DB
	username := c.MustGet("user").(string)

	res, err := db.GetMediaRequestList(model.User{Username: username})
	if err != nil {
		if db.IsMediaItemNotFound(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Error("Server: media request list get failed: %+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) handleAddMediaRequest(c *gin.Context) {
	db := s.config.DB
	username := c.MustGet("user").(string)

	var in model.MediaRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		log.Debug("Server: media request add failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := db.GetUser(model.User{Username: username})
	if err != nil {
		log.Error("Server: media request add failed: +v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	in.User = user

	if err := db.AddMediaRequest(in); err != nil {
		log.Error("Server: media request add failed: %+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (s *Server) handleGetMediaResponseList(c *gin.Context) {
	db := s.config.DB
	username := c.MustGet("user").(string)

	res, err := db.GetMediaResponseList(model.User{Username: username})
	if err != nil {
		if db.IsMediaItemNotFound(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Error("Server: media response list get failed: %+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) handleAddMediaResponse(c *gin.Context) {
	db := s.config.DB
	username := c.MustGet("user").(string)

	var in model.MediaResponse
	if err := c.ShouldBindJSON(&in); err != nil {
		log.Debug("Server: media response add failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	owner, err := db.GetUser(model.User{Username: username})
	if err != nil {
		log.Error("Server: media response add failed: +v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	in.Owner = owner

	if err := db.AddMediaResponse(in); err != nil {
		log.Error("Server: media response add failed: %+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, nil)
}