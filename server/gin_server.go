package server

import (
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
)

// GinServer gin server inherits
type GinServer struct {
	Server
	*gin.Engine // Engine *gin.Engine
}

// NewGinServer create a gin web server
func NewGinServer(env string, config ...*simplejson.Json) (*GinServer, error) {
	server := &GinServer{}
	server.Env = env
	var err error
	server.AppPath, server.Config, err = initialize(env, config...)
	if err != nil {
		return nil, err
	}

	if env == "pro" {
		gin.SetMode(gin.ReleaseMode)
		server.Engine = gin.New()
	} else {
		server.Engine = gin.Default()
	}

	return server, nil
}

// Run start GinServer run
func (s *GinServer) Run(addr ...string) error {
	var address string
	if len(addr) == 0 {
		address = s.Config.Get("app").Get("addr").MustString()
	} else {
		address = addr[0]
	}
	err := s.Engine.Run(address)
	if err != nil {
		return err
	}

	return nil
}
