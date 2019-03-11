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

	var in struct {
		Username  string  `json:"username" form:"username" binding:"required"`
		Password  *string `json:"password,omitempty" form:"password" binding:"required"`
		Password1 *string `json:"password1,omitempty" form:"password1"`
		Token     string  `json:"token"`
	}

	if err := c.ShouldBindJSON(&in); err != nil {
		log.Debug("Server: user creation failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := model.User{
		Username: in.Username,
	}

	if _, err := db.GetUser(user); err == nil {
		log.Debug("Server: user creation failed: user already exists: %s", in.Username)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	if *in.Password != *in.Password1 {
		log.Debug("Server: user creation failed: passwords mismatch: %s - %s", in.Password, in.Password1)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ph, err := bcrypt.GenerateFromPassword([]byte(*in.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Server: user creation failed: %+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	user.PassHash = string(ph)

	if err := db.CreateUser(user); err != nil {
		log.Error("Server: user creation failed: %+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := s.createJwtToken(user)
	if err != nil {
		log.Error("Server: user creation failed: %+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	in.Password = nil
	in.Password1 = nil
	in.Token = token

	c.JSON(http.StatusCreated, in)
}

func (s *Server) handleLoginUser(c *gin.Context) {
	var in struct {
		Username string  `json:"username" form:"username" binding:"required"`
		Password *string `json:"password,omitempty" form:"password" binding:"required"`
		Token    string  `json:"token"`
	}

	if err := c.ShouldBindJSON(&in); err != nil {
		log.Debug("Server: user login failed: %+v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := s.config.DB.GetUser(model.User{Username: in.Username})
	if err != nil {
		log.Debug("Server: user login failed: %+v", err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(*in.Password)); err != nil {
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

	in.Password = nil
	in.Token = token

	c.JSON(http.StatusOK, in)
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
