package main

import (
	"fmt"
<<<<<<< HEAD
=======
	"github.com/A-walker-ninght/miniRedis/resp/handler"
>>>>>>> 70f3717 (resp 2023.3.1)
	"os"

	"github.com/A-walker-ninght/miniRedis/config"
	"github.com/A-walker-ninght/miniRedis/lib/logger"
	"github.com/A-walker-ninght/miniRedis/tcp"
)

const configFile string = "redis.conf"

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6380,
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})
	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}
	err := tcp.ListenAndServeWithSignal(
		&tcp.Config{
			Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
		},
<<<<<<< HEAD
		tcp.MakeClientPool(),
=======
		handler.MakeHandler(),
>>>>>>> 70f3717 (resp 2023.3.1)
	)
	if err != nil {
		logger.Error(err)
	}
}
