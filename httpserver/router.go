package httpserver

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/longweb/httpserver/handlers"
)

func InitRoute(dotserver *dotweb.Dotweb) {
	dotserver.HttpServer.GET("/", handlers.Index)
	dotserver.HttpServer.GET("/mstate", handlers.Memstate)
	dotserver.HttpServer.GET("/testauth", handlers.TestAuth)
	dotserver.HttpServer.GET("/testmessage", handlers.TestMessage)

	dotserver.HttpServer.GET("/state", handlers.State)
	dotserver.HttpServer.POST("/sendmessage", handlers.SendMessage)
	dotserver.HttpServer.ServerFile("/static/*filepath", "/home/emoney/longweb/www")
	dotserver.HttpServer.ServerFile("/js/*filepath", "/home/emoney/longweb/javascript")
	dotserver.HttpServer.WebSocket("/ws/onsocket", handlers.OnWebSocket)
	dotserver.HttpServer.HiJack("/poll/onpolling", handlers.OnPolling)
}
