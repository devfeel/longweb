# LongWeb
Simple and easy go realtime-web gateway

[![GoDoc](https://godoc.org/github.com/devfeel/longweb?status.svg)](https://godoc.org/github.com/devfeel/longweb)
[![Go Report Card](https://goreportcard.com/badge/github.com/devfeel/longweb)](https://goreportcard.com/report/github.com/devfeel/longweb)

## 1. Install

```
go get -u github.com/devfeel/longweb
```

## 2. Features
* 支持Websocket\longpoll，消灭浏览器兼容之痛
* 原有业务系统无痛接入
* 支持公开与授权模式(token)
* 支持配置化部署
* 支持连接数据持久化，目前支持influxdb

 
## 3. 主要文件说明
* httpseerver/handlers/wshandler  处理websocket请求
* httpseerver/handlers/pollhandler 处理长轮询请求（用于不支持websocket的访问端）
* httpseerver/handlers/apihandler 处理与后端业务系统的交互
 
* message/message 后端业务系统的消息接收、处理、推送
* message/userclient 用户连接的抽象
* message/usergroup 用户组别    
* 其中三层：app-》group-》client
 
## 4. 应用场景
* 原有常规http系统需要支持websocket等长连接功能，无需开发改造，对接该网关轻松实现
* 实时在线用户数统计系统
* 聊天室类应用
* 用户行为实时监测及反馈系统
