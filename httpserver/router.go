package httpserver

import (
	"devfeel/dotweb"
	"devfeel/longweb/httpserver/handlers"
)

func InitRoute(dotweb *dotweb.Dotweb) {
	dotweb.HttpServer.GET("/", handlers.Index)
	dotweb.HttpServer.GET("/mstate", handlers.Memstate)
	dotweb.HttpServer.GET("/testauth", handlers.TestAuth)
	dotweb.HttpServer.GET("/testmessage", handlers.TestMessage)

	dotweb.HttpServer.GET("/state", handlers.State)
	dotweb.HttpServer.POST("/sendmessage", handlers.SendMessage)
	dotweb.HttpServer.ServerFile("/static/*filepath", "/home/emoney/longweb/www")
	dotweb.HttpServer.ServerFile("/js/*filepath", "/home/emoney/longweb/javascript")
	dotweb.HttpServer.WebSocket("/ws/onsocket", handlers.OnWebSocket)
	dotweb.HttpServer.HiJack("/poll/onpolling", handlers.OnPolling)
}
