package server

import (
	"fmt"
	"github.com/cnaize/mz-center/server/handle"
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

	r.GET("/", handle.Status)
	v1 := r.Group("/v1", handle.Log)
	{
		searches := v1.Group("/searches")
		{
			searches.GET("", handle.ListSearches)
			searches.GET("/:text", handle.GetSearch)
			searches.POST("/:text", handle.AddSearch)
			searches.POST("/:text/:user", handle.AddSearchResult)
		}
	}

	return &Server{
		config: config,
		router: r,
	}
}

func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf(":%d", s.config.Port))
}
