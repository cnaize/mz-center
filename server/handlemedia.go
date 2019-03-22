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

	var inRequest model.MediaRequest
	if err := c.ShouldBindJSON(&inRequest); err != nil || inRequest.Owner.Username == "" {
		log.Debug("Server: media request add failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// TODO:
	//  handle "protected" mode

	if inRequest.Mode != model.MediaAccessTypePublic && inRequest.Mode != model.MediaAccessTypePrivate {
		inRequest.Mode = model.MediaAccessTypePublic
	}
	inRequest.User = model.User{Username: username}

	if err := db.AddMediaRequest(inRequest); err != nil {
		log.Error("Server: media request add failed: %+v", err)
		if db.IsMediaItemNotFound(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

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

	var inResponse model.MediaResponse
	if err := c.ShouldBindJSON(&inResponse); err != nil || inResponse.MediaRequest.User.Username == "" {
		log.Debug("Server: media response add failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// TODO:
	//  handle "protected" mode

	if inResponse.MediaRequest.Mode != model.MediaAccessTypePublic && inResponse.MediaRequest.Mode != model.MediaAccessTypePrivate {
		inResponse.MediaRequest.Mode = model.MediaAccessTypePublic
	}
	inResponse.MediaRequest.Owner = model.User{Username: username}

	if err := db.AddMediaResponse(inResponse); err != nil {
		log.Error("Server: media response add failed: %+v", err)
		if db.IsMediaItemNotFound(err) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, nil)
}
