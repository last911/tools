package server

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/last911/tools"
	"github.com/last911/tools/log"
	"io/ioutil"
)

// Server web server
type Server struct {
	Config  *simplejson.Json // Config is Conf struct
	AppPath string           // AppPath bin path
	Env     string           // Environment
}

// initialize return app path and config json or error
func initialize(env string, conf ...*simplejson.Json) (string, *simplejson.Json, error) {
	var err error
	appPath, err := tools.AbsolutePath()
	if err != nil {
		return "", nil, err
	}
	var config *simplejson.Json
	if len(conf) > 0 {
		config = conf[0]
	} else {
		configPath := appPath + "conf/config-" + env + ".json"
		log.Debug(fmt.Sprintf("load config from[%s]", configPath))

		b, err := ioutil.ReadFile(configPath)
		if err != nil {
			return "", nil, err
		}

		config, err = simplejson.NewJson(b)
		if err != nil {
			return "", nil, err
		}
	}

	logConf := config.Get("log")
	logConfPath, err := logConf.String()
	var opts = ""
	if err == nil {
		b, err := ioutil.ReadFile(logConfPath)
		if err != nil {
			return "", nil, err
		}
		opts = string(b)
	}
	log.InitLogger(opts)

	return appPath, config, nil
}
