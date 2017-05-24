package message

import (
	"github.com/devfeel/dotweb"

	"errors"
	"fmt"
	. "github.com/devfeel/longweb/const"
	"github.com/devfeel/longweb/exception"
	"github.com/devfeel/longweb/framework/json"
	"github.com/devfeel/longweb/framework/log"
	"sync"
	"sync/atomic"
	"time"
)

var clientIndex uint64
var clientCreateCount uint64

const DefaultTimeOut = 60 //second

func GetMaxClientIndex() uint64 {
	return clientIndex
}

func GetTotalClientCreateCount() uint64 {
	return clientCreateCount
}

//应用集合
var AppPool map[string]*AppGroups
var appLock *sync.RWMutex

//客户端池
var clientPool sync.Pool

//定时器池 -- *Timer
var timerPool sync.Pool

func init() {
	AppPool = make(map[string]*AppGroups)
	appLock = new(sync.RWMutex)
	clientPool = sync.Pool{
		New: func() interface{} {
			atomic.AddUint64(&clientCreateCount, 1)
			return &UserClient{}
		},
	}
	timerPool = sync.Pool{
		New: func() interface{} {
			return time.NewTimer(DefaultTimeOut * time.Second)
		},
	}
}

//user client for websocket
type UserClient struct {
	Index        uint64 //客户端唯一索引
	ConnType     string
	webSocket    *dotweb.WebSocket
	httpContext  dotweb.Context
	MessageChan  chan string `json:"-"`
	isHijackSend bool
	timer        *time.Timer
	TimeOut      int64  //超时时间 - 单位为秒
	From         string //请求来源，比如site、mobile、jrpt等
	UserID       string //user id
	GroupId      string //用户组编号
	AppId        string //用户应用编号
	RemoteIP     string //用户IP信息“ip:port”
	ReferrerUrl  string //ReferrerUrl
	IsAuth       bool   //是否带鉴权
}

//create a new UserClient with socketconn&userinfo
func NewClient(appId, userId, groupId, from string, isAuth bool, ws *dotweb.WebSocket, context dotweb.Context) *UserClient {
	atomic.AddUint64(&clientIndex, 1)
	client := clientPool.Get().(*UserClient)
	client.Reset(appId, userId, groupId, from, isAuth, ws, context)
	//log the new client info
	logger.Log("UserClient:NewNormalClient["+client.GetClientInfo()+"] Connect", LogTarget_UserClient, LogLevel_Debug)
	return client
}

//register new userclient with roomid
//return values:
//-10001: not exists appid
//-10002: usergroup error
//0: ok
func RegisterClient(client *UserClient) (*UserGroup, int) {
	app, exists := GetAppGroups(client.AppId)
	if !exists {
		return nil, -10001
	}
	userGroup := app.GetAndInitUserGroup(client.AppId, client.GroupId)
	if userGroup == nil {
		return nil, -10002
	}
	userGroup.AddClient(client)
	return userGroup, 0
}

//remove userclient with roomid
func RemoveClient(client *UserClient) {
	if client == nil {
		logString := "UserClient::RemoveClient -> [client is nil]"
		logger.Log(logString, LogTarget_UserClient, LogLevel_Debug)
		return
	}
	group, exists := GetUserGroup(client.AppId, client.GroupId)
	if exists {
		group.DeleteClient(client)
	}
	if client.webSocket != nil {
		client.webSocket.Conn.Close()
	}
	if client.httpContext != nil {
		//检查是否hijack模式
		if client.httpContext.IsHijack() {
			client.httpContext.HijackConn().Close()
		}
		if client.timer != nil {
			//归还timer对象
			client.timer.Stop()
			timerPool.Put(client.timer)
		}
	}

	//记录访问日志
	logString := "UserClient::RemoveClient -> [" + client.GetClientInfo() + "]"
	logger.Log(logString, LogTarget_UserClient, LogLevel_Debug)

	client.Reset("", "", "", "", false, nil, nil)
	clientPool.Put(client)
}

//reset userclient attr
func (uc *UserClient) Reset(appId, userId, groupId, from string, isAuth bool, ws *dotweb.WebSocket, context dotweb.Context) {
	uc.AppId = appId
	uc.UserID = userId
	uc.GroupId = groupId
	uc.From = from
	uc.webSocket = ws
	uc.httpContext = context
	uc.IsAuth = isAuth
	uc.MessageChan = nil
	uc.timer = nil
	uc.isHijackSend = false

	uc.ConnType = ""
	uc.RemoteIP = ""
	uc.ReferrerUrl = ""
	if uc.webSocket != nil {
		uc.ConnType = ConnType_WebSocket
	} else if uc.httpContext != nil {
		uc.MessageChan = make(chan string, 1)
		uc.ConnType = ConnType_LongPoll
	}
	if uc.webSocket == nil && uc.httpContext == nil {
		uc.RemoteIP = ""
		uc.ReferrerUrl = ""
	} else {
		uc.RemoteIP = uc.GetRemoteAddr()
		uc.ReferrerUrl = uc.GetReferrerUrl()
	}
}

//send message to client
func (uc *UserClient) SendMessage(message string) {
	defer func() {
		if err := recover(); err != nil {
			ex := exception.CatchError("UserClient::SendMessage", err)
			//记录访问日志
			logString := "UserClient::SendMessage -> to:[" + fmt.Sprint(uc) + "] send:[" + message + "] error:[" + ex.GetErrString() + "]"
			logger.Log(logString, LogTarget_UserClient, LogLevel_Error)
		}
	}()
	if uc.ConnType == ConnType_WebSocket {
		uc.webSocket.SendMessage(message)
	} else if uc.ConnType == ConnType_LongPoll {
		if uc.MessageChan != nil {
			//置入消息
			uc.MessageChan <- message
			uc.isHijackSend = true
		}
	}
}

//get client's remoteaddr
func (uc *UserClient) GetRemoteAddr() string {
	if uc.ConnType == ConnType_WebSocket {
		return uc.webSocket.Conn.Request().RemoteAddr
	} else if uc.ConnType == ConnType_LongPoll {
		return uc.httpContext.RemoteIP()
	}
	return ""
}

//read message from websocket.conn or hijackConn
func (uc *UserClient) ReadMessage() (string, error) {
	if uc.ConnType == ConnType_WebSocket {
		return uc.webSocket.ReadMessage()
	} else if uc.ConnType == ConnType_LongPoll {
		//重置定时器
		uc.timer = timerPool.Get().(*time.Timer)
		uc.timer.Reset(time.Second * time.Duration(uc.TimeOut))
		for {
			select {
			case str := <-uc.MessageChan:
				return str, nil
			case <-uc.timer.C:
				{
					return "", errors.New("time out")
				}
			}
		}
	} else {
		return "", nil
	}
}

//get client's websocket
func (uc *UserClient) GetWebSocket() *dotweb.WebSocket {
	return uc.webSocket
}

//get client's referrer url
func (uc *UserClient) GetReferrerUrl() string {
	if uc.ConnType == ConnType_WebSocket {
		return uc.webSocket.Conn.Request().Referer()
	}
	if uc.ConnType == ConnType_LongPoll {
		return uc.httpContext.Request().Referer()
	}
	return ""
}

//get client's json string
func (uc *UserClient) GetClientInfo() string {
	//return jsonutil.GetJsonString(uc)
	str, err := jsonutil.Marshal(uc)
	if err != nil {
		logger.Error("UserClient:GetClientInfo error => "+err.Error(), LogTarget_UserClient)
	}
	return str
}

//get client's serverindex
func (uc *UserClient) GetIndex() uint64 {
	return uc.Index
}
