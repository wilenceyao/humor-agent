package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-agent/config"
)

var DefaultAdminServer *Server

type Server struct {
	router       *gin.Engine
	adminService *Service
}

func RunAdminServer() error {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.Logger
	DefaultAdminServer = &Server{
		router:       gin.Default(),
		adminService: &Service{},
	}

	DefaultAdminServer.addApi()
	return DefaultAdminServer.router.Run(fmt.Sprintf(":%d", config.Config.Admin.Port))
}

func (s *Server) addApi() {
	//s.router.POST("tts", s.adminService.Tts)
}
