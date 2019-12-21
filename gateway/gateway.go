package main

import (
	"log"
	"os"

	"github.com/zenchan/server-go/common/config"
	"github.com/zenchan/server-go/common/utils"
	"github.com/zenchan/server-go/libs/xlog"
)

func main() {
	utils.FlagParse()

	if err := config.Load(utils.ConfigFile); err != nil {
		log.Printf("load config failed: %s\n", err.Error())
		os.Exit(-1)
	}

	pname := utils.ProcessName()

	if err := xlog.InitLogging(config.SrvCfg.LogPath+"/"+pname, config.SrvCfg.LogLevel); err != nil {
		log.Printf("init logger failed: %s\n", err.Error())
		os.Exit(-1)
	}

	xlog.Info(pname + " run successfully")
}
