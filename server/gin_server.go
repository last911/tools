package server

import (
	"github.com/gin-gonic/gin"
)

// GinServer gin server inherits
type GinServer struct {
	Server
	Engine *gin.Engine // Engine *gin.Engine
}

// NewGinServer create a gin web server
func NewGinServer(env string, configStr ...string) (*GinServer, error) {
	server := &GinServer{}
	server.Env = env
	var err error
	if len(configStr) > 0 {
		server.AppPath, server.Config, err = initializeFromConfig(configStr[0])
	} else {
		server.AppPath, server.Config, err = initialize(env)
	}
	if err != nil {
		return nil, err
	}

	if env == "pro" {
		server.Engine = gin.New()
	} else {
		server.Engine = gin.Default()
	}

	return server, nil
}

// Run start GinServer run
func (s *GinServer) Run(addr ...string) error {
	if len(addr) == 0 {
		addr[0] = s.Config.Get("app").Get("addr").MustString()
	}
	err := s.Engine.Run(addr[0])
	if err != nil {
		return err
	}

	return nil
}
