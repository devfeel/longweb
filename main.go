/*devfeel.longweb
* Author: Panxinming
* LastUpdateTime: 2017-02-07 10:00
 */
package main

import (
	"flag"
	"fmt"
	"github.com/devfeel/longweb/config"
	"github.com/devfeel/longweb/framework/file"
	"github.com/devfeel/longweb/framework/log"
	"github.com/devfeel/longweb/httpserver"
	"github.com/devfeel/longweb/message"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	innerLogger *logger.InnerLogger
	configFile  string
)

func init() {
	innerLogger = logger.GetInnerLogger()
}
func main() {
	defer func() {
		if err := recover(); err != nil {
			strLog := "longweb:main recover error => " + fmt.Sprintln(err)
			os.Stdout.Write([]byte(strLog))
			innerLogger.Error(strLog)

			buf := make([]byte, 4096)
			n := runtime.Stack(buf, true)
			innerLogger.Error(string(buf[:n]))
			os.Stdout.Write(buf[:n])
		}
	}()

	currentBaseDir := fileutil.GetCurrentDirectory()
	flag.StringVar(&configFile, "config", "", "配置文件路径")
	if configFile == "" {
		configFile = currentBaseDir + "/app.conf"
	}
	//启动内部日志服务
	logger.StartInnerLogHandler(currentBaseDir)

	//加载xml配置文件
	appconfig := config.InitConfig(configFile)

	//设置基本目录
	config.SetBaseDir(currentBaseDir)

	//启动日志服务
	logger.StartLogHandler(appconfig.Log.FilePath)

	//初始化应用信息
	message.InitAppInfo()

	//start message service
	message.StartMessageService()

	//监听系统信号
	go listenSignal()

	err := httpserver.StartServer()
	if err != nil {
		innerLogger.Warn("HttpServer.StartServer失败 " + err.Error())
		fmt.Println("HttpServer.StartServer失败 " + err.Error())
	}

}

func listenSignal() {
	c := make(chan os.Signal, 1)
	//syscall.SIGSTOP
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		innerLogger.Info("signal::ListenSignal [" + s.String() + "]")
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP: //配置重载
			innerLogger.Info("signal::ListenSignal reload config begin...")
			//重新加载配置文件
			config.InitConfig(configFile)
			//初始化应用信息
			message.InitAppInfo()
			//start message service
			message.ReStartMessageService()
			innerLogger.Info("signal::ListenSignal reload config end")
		default:
			return
		}
	}
}
