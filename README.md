# longweb
为web程序提供长连接网关服务
 
##主要文件说明
* httpseerver/handlers/wshandler  处理websocket请求
* httpseerver/handlers/pollhandler 处理长轮询请求（用于不支持websocket的访问端）
* httpseerver/handlers/apihandler 处理与后端业务系统的交互
 
* message/message 后端业务系统的消息接收、处理、推送
* message/userclient 用户连接的抽象
* message/usergroup 用户组别    
* 其中三层：app-》group-》client
