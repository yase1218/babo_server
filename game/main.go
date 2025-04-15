package main

import (
	"babo/game/app"
	_ "babo/game/handler"
	"babo/utility/application"
	"babo/utility/zlog"
	"flag"

	"go.uber.org/zap"
)

const (
	APP_NAME = "gameserver"
)

var (
	confFile    = flag.String("conf", "./config/config.yml", "config file path")
	cancelPrint = flag.Bool("cancelprint", true, "print log to console")
	closeDebug  = flag.Bool("closedebug", true, "close debug module")
	g_app       app.GameServerApp
)

func main() {
	flag.Parse()
	if err := application.StartApp(&g_app, APP_NAME, *confFile, g_app.GetConfig(), *cancelPrint, false); err != nil {
		zlog.Error("Start failed", zap.String("err", err.Error()))
		return
	}
}
