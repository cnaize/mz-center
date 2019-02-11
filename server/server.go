package server

import (
	"fmt"
	"github.com/cnaize/mz-common/log"
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
	v1 := r.Group("/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/signup", s.handleCreateUser)
			users.POST("/signin", s.handleLoginUser)
		}
		
		searches := v1.Group("/searches", s.handleAuthCheck)
		{
			reqs := searches.Group("/requests")
			{
				reqs.GET("", s.handleGetSearchRequestList)
				reqs.POST("", s.handleAddSearchRequest)
			}
			resps := searches.Group("/responses")
			{
				resps.GET("", s.handleGetSearchResponseList)
				resps.POST("", s.handleAddSearchResponseList)
			}
		}
	}

	return s
}

func (s *Server) Run() error {
	log.Info("MuzeZone Center: running server on port: %d", s.config.Port)
	return s.router.Run(fmt.Sprintf(":%d", s.config.Port))
}
