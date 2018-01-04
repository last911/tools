package server

import (
	"github.com/last911/tools"
	"fmt"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
)

type Server struct {
	Config  *simplejson.Json // Config is Conf struct
	AppPath string           // AppPath bin path
	Env     string           // Environment
}

func initialize(env string) (appPath string, config *simplejson.Json, err error) {
	appPath, err = tools.AbsolutePath()
	if err != nil {
		return
	}
	configPath := appPath + "conf/config-" + env + ".json"
	log.Debug(fmt.Sprintf("load config from[%s]", configPath))

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}

	return initializeFromConfig(string(b))
}

func initializeFromConfig(configStr string) (appPath string, config *simplejson.Json, err error) {
	config, err = simplejson.NewJson([]byte(configStr))
	if err != nil {
		return
	}

	log.Debug("config data:", config)
	logConf := config.Get("log")
	logConfPath, err := logConf.String()
	var logger log.LoggerInterface
	if err != nil {
		logger = log.Default
		err = nil
	} else {
		logger, err = log.LoggerFromConfigAsFile(logConfPath)
		if err != nil {
			return
		}
	}
	log.UseLogger(logger)

	return
}
