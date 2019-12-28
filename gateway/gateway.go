package main

import (
	"fmt"
	"log"
	"net"
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
	if err := xlog.InitLogging(config.LogPath()+"/"+pname, config.LogLevel(),
		xlog.WithStdout(true),
	); err != nil {
		log.Printf("init logger failed: %s\n", err.Error())
		os.Exit(-1)
	}

	// if err := listen(); err != nil {
	// 	os.Exit(-1)
	// }

	xlog.Info(pname + " run successfully")
	xlog.SetStdout(false)
}

func listen() (err error) {
	var (
		tcpLis net.Listener
	)

	if config.TCPPort() != 0 {
		addr := fmt.Sprintf("%d", config.TCPPort())
		if tcpLis, err = net.Listen("tcp", addr); err != nil {
			xlog.Infof("listen tcp %d failed: %s", config.TCPPort(), err.Error())
			return
		}
		_ = tcpLis
	}
	return
}
