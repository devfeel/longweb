package httpserver

import (
	"devfeel/longweb/config"
	"devfeel/longweb/framework/log"
	"devfeel/longweb/message"
	"github.com/devfeel/dotweb"
	"strconv"
)

func StartServer() error {

	//初始化DotServer
	dotserver := dotweb.New()

	//设置dotserver日志目录
	dotserver.SetLogPath(config.CurrentConfig.Log.FilePath)

	//设置路由
	InitRoute(dotserver)

	innerLogger := logger.GetInnerLogger()

	//启动监控服务
	pprofport := config.CurrentConfig.HttpServer.PProfPort
	go dotserver.StartPProfServer(pprofport)

	// 开始服务
	port := config.CurrentConfig.HttpServer.HttpPort
	innerLogger.Debug("dotweb.StartServer => " + strconv.Itoa(port))
	err := dotserver.StartServer(port)
	return err
}

func ReSetServer() {
	//初始化应用信息
	message.InitAppInfo()
}
