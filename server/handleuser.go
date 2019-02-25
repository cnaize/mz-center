package server

import (
	"github.com/cnaize/mz-common/log"
	"github.com/cnaize/mz-common/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (s *Server) handleCreateUser(c *gin.Context) {
	db := s.config.DB

	var inUser model.User
	if err := c.ShouldBindJSON(&inUser); err != nil {
		log.Debug("Server: user creation failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if _, err := db.GetUser(inUser); err == nil {
		log.Debug("Server: user creation failed: user already exists: %s", inUser.Username)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	if *inUser.Password != *inUser.Password1 {
		log.Debug("Server: user creation failed: passwords mismatch: %s - %s", *inUser.Password, *inUser.Password1)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ph, err := bcrypt.GenerateFromPassword([]byte(*inUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Server: user creation failed: %+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	inUser.PassHash = string(ph)

	if err := db.CreateUser(inUser); err != nil {
		log.Error("Server: user creation failed: %+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := s.createJwtToken(inUser)
	if err != nil {
		log.Error("Server: user creation failed: %+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	inUser.Token = token

	c.JSON(http.StatusCreated, inUser)
}

func (s *Server) handleLoginUser(c *gin.Context) {
	var inUser model.User
	if err := c.ShouldBindJSON(&inUser); err != nil {
		log.Debug("Server: user login failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := s.config.DB.GetUser(inUser)
	if err != nil {
		log.Debug("Server: user login failed: %+v", err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(*inUser.Password)); err != nil {
		log.Debug("Server: user login failed: %+v", err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	token, err := s.createJwtToken(user)
	if err != nil {
		log.Warn("Server: user login failed: %+v", err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	user.Token = token

	c.JSON(http.StatusOK, user)
}

func (s *Server) createJwtToken(user model.User) (string, error) {
	tk := model.Token{Username: user.Username}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tk)
	tokenStr, err := token.SignedString([]byte(s.config.JwtTokenPassword))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}