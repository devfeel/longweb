package httpserver

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/longweb/httpserver/handlers"
)

func InitRoute(dotserver *dotweb.DotWeb) {
	dotserver.HttpServer.Router().GET("/", handlers.Index)
	dotserver.HttpServer.Router().GET("/mstate", handlers.Memstate)
	dotserver.HttpServer.Router().GET("/testauth", handlers.TestAuth)
	dotserver.HttpServer.Router().GET("/testmessage", handlers.TestMessage)

	dotserver.HttpServer.Router().GET("/state", handlers.State)
	dotserver.HttpServer.Router().GET("/statedata", handlers.StateData)
	dotserver.HttpServer.Router().POST("/sendmessage", handlers.SendMessage)
	dotserver.HttpServer.Router().ServerFile("/www/*filepath", "/home/emoney/longweb/www")
	dotserver.HttpServer.Router().WebSocket("/ws/onsocket", handlers.OnWebSocket)
	dotserver.HttpServer.Router().HiJack("/poll/onpolling", handlers.OnPolling)
}
