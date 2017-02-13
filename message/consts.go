package message

//连接类型
const (
	ConnType_WebSocket = "websocket"
	ConnType_LongPoll  = "longpoll"
)

//消息级别
const (
	MessageLevel_All    = "0" //"all"
	MessageLevel_Normal = "1" //"normal"
	MessageLevel_Auth   = "2" //"auth" 表示仅发送给经过鉴权的连接
)
