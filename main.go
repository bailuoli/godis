package main

import (
	"fmt"
	"godis/config"
	"godis/lib/logger"
	"godis/tcp"
	"os"
)

var banner = `
   ______          ___
  / ____/___  ____/ (_)____
 / / __/ __ \/ __  / / ___/
/ /_/ / /_/ / /_/ / (__  )
\____/\____/\__,_/_/____/
`

// 配置文件名称
const configFileName = "redis.conf"

// 没有配置文件则使用该配置
var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6378,
}

// 判断文件是否存在
func fileExists(fileName string) bool {
	stat, err := os.Stat(fileName)
	return err != nil && !stat.IsDir()
}

func main() {
	println(banner)
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	//如果配置文件存在，则使用配置文件，如果配置文件不存在，则使用默认配置
	if fileExists(configFileName) {
		config.SetupConfig(configFileName)
	} else {
		config.Properties = defaultProperties
	}
	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	},
		tcp.MakeEchoHandler())
	if err != nil {
		logger.Error(err)
	}

}
