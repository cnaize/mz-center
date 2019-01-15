package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config Config
	router *gin.Engine
}

func New(config Config) *Server {
	r := gin.Default()
	r.Use(cors.Default())

	s := &Server{
		config: config,
		router: r,
	}

	r.GET("/", s.handleStatus)
	v1 := r.Group("/v1", s.handleLog)
	{
		searches := v1.Group("/searches")
		{
			searches.GET("", s.handleGetSearchRequestList)
			searches.POST("/:text", s.handleAddSearchRequest)
			searches.GET("/:text", s.handleGetSearchResponseList)
			searches.POST("/:text/:username", s.handleAddSearchResponseList)
		}
	}

	return s
}

func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf(":%d", s.config.Port))
}
