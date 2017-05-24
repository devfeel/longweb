package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devfeel/longweb/config"
	. "github.com/devfeel/longweb/const"
	"github.com/devfeel/longweb/exception"
	"github.com/devfeel/longweb/framework/log"
	"github.com/devfeel/longweb/framework/task"
	"strconv"
	"sync"
)

var msgQueueMap map[string]*MessageQueue
var msgLock *sync.RWMutex

const MaxQueueMessageCount = 10000

func init() {
	//max message count with wait
	msgQueueMap = make(map[string]*MessageQueue)
	msgLock = new(sync.RWMutex)
}

//重新启动消息监听与消费服务
func ReStartMessageService() {
	task.StopAllTask()

	apps := config.GetAppList()
	for _, app := range apps {
		//init msgqueue
		msgLock.Lock()
		msgQueueMap[app.AppID] = &MessageQueue{
			AppID:   app.AppID,
			MsgChan: make(chan Message, MaxQueueMessageCount),
		}
		msgLock.Unlock()

		//init msg task
		t := task.NewTask(app.AppID, task_DealMessage)
		t.TaskData = app.AppID
		task.RemoveTask(t.Name)
		task.RegisterTask(t)

		logger.Debug("ReStartMessageService => "+fmt.Sprintln(app), LogTarget_Message)
	}

	//启动所有消息任务
	task.StartAllTask()
}

//全局启动消息监听与消费服务
func StartMessageService() {
	logger.Debug("StartMessageService => begin", LogTarget_Message)

	for _, app := range config.GetAppList() {

		msgLock.Lock()
		//init msgqueue
		msgQueueMap[app.AppID] = &MessageQueue{
			AppID:   app.AppID,
			MsgChan: make(chan Message, MaxQueueMessageCount),
		}
		msgLock.Unlock()

		//init msg task
		t := task.NewTask(app.AppID, task_DealMessage)
		t.TaskData = app.AppID
		task.RegisterTask(t)

		logger.Debug("StartMessageService => "+fmt.Sprintln(app), LogTarget_Message)
	}

	//启动所有消息任务
	task.StartAllTask()

	logger.Debug("StartMessageService => end", LogTarget_Message)
}

//从队列读取消息并处理
func task_DealMessage(task *task.TaskInfo) {
	defer func() {
		if err := recover(); err != nil {
			ex := exception.CatchError("message::task_DealMessage", err)
			//记录访问日志
			logString := "message::task_DealMessage[" + fmt.Sprint(task.TaskData) + "] error:[" + ex.GetErrString() + "]"
			logger.Log(logString, LogTarget_Message, LogLevel_Error)
		}
	}()

	appId, isOk := task.TaskData.(string)
	if !isOk {
		logger.Error("DealMessage["+fmt.Sprintln(task)+"] error => taskdata can't convert to string", LogTarget_Message)
		return
	}
	msg, err := readMessage(appId)

	if err != nil {
		logger.Error("DealMessage["+fmt.Sprintln(task)+"] error => ["+err.Error()+"]", LogTarget_Message)
		return
	}

	if msg.AppID == "" {
		logger.Error("DealMessage["+fmt.Sprintln(task)+"] error => AppID is empty ["+fmt.Sprintln(msg)+"]", LogTarget_Message)
		return
	}
	if msg.ToAppID == "" {
		logger.Error("DealMessage["+fmt.Sprintln(task)+"] error => ToAppID is empty ["+fmt.Sprintln(msg)+"]", LogTarget_Message)
		return
	}

	//暂时屏蔽跨App转发
	if msg.ToAppID != msg.AppID {
		logger.Error("DealMessage["+fmt.Sprintln(task)+"] error => AppID not equal ToAppID ["+fmt.Sprintln(msg)+"]", LogTarget_Message)
		return
	}

	var app *AppGroups
	var group *UserGroup
	var client *UserClient
	var exists bool

	//获取应用用户群组
	app, exists = GetAppGroups(msg.ToAppID)
	if !exists {
		logger.Error("DealMessage["+fmt.Sprintln(task)+"]:GetAppGroups error => not exists ToAppID ["+fmt.Sprintln(msg)+"]", LogTarget_Message)
		return
	}

	//*********发送逻辑**********

	//发送给整个应用用户群
	//tips:暂未过滤
	if msg.ToGroupID == "" && msg.ToUserID == "" {
		logger.Debug("DealMessage["+fmt.Sprintln(task)+"]: Begin SendMessage [ToApp] ["+fmt.Sprintln(msg)+"]", LogTarget_Message)
		count := app.SendMessage(msg)
		logger.Debug("DealMessage["+fmt.Sprintln(task)+"]: End SendMessage [ClientCount="+strconv.Itoa(count)+"] [ToApp] ["+fmt.Sprintln(msg)+"]", LogTarget_Message)
		//删除标记需要删除的Client
		return
	}

	//获取用户组信息
	if msg.ToGroupID != "" {
		group, exists = app.GetUserGroup(msg.ToGroupID)
		if !exists {
			logger.Error("DealMessage["+fmt.Sprintln(task)+"]:GetUserGroup error => not exists ToGroupID ["+fmt.Sprintln(msg)+"]", LogTarget_Message)
			return
		}
	}

	//发送给用户组
	if msg.ToUserID == "" {
		logger.Debug("DealMessage["+fmt.Sprintln(task)+"]: Begin SendMessage [ToGroup] ["+fmt.Sprintln(msg)+"]", LogTarget_Message)
		count := group.SendMessage(msg)
		logger.Debug("DealMessage["+fmt.Sprintln(task)+"]: End SendMessage [ClientCount="+strconv.Itoa(count)+"] [ToGroup] ["+fmt.Sprintln(msg)+"]", LogTarget_Message)
		return
	}

	//发送给单用户
	if msg.ToUserID != "" {
		client, exists = group.GetUserClient(msg.ToUserID)
		if !exists {
			logger.Error("DealMessage["+fmt.Sprintln(task)+"]:GetUserClient error => not exists ToUserID ["+fmt.Sprintln(msg)+"]", LogTarget_Message)
			return
		} else {
			logger.Debug("DealMessage["+fmt.Sprintln(task)+"]:SendMessage [ToUser] ["+fmt.Sprintln(msg)+"]", LogTarget_Message)
			client.SendMessage(msg.Content)
		}

	}
}

//************消息相关******************

//消息队列
type MessageQueue struct {
	AppID   string
	MsgChan chan Message
}

//消息定义
type Message struct {
	//消息关联AppID
	AppID string
	//消息关联GroupID
	GroupID string
	//消息关联UserID
	UserID string
	//消息接收级别 - 目前分为all、auth、normal -匹配MessageLevel_Const
	MessageLevel string
	//消息接收AppID
	ToAppID string
	//消息接收GroupID
	ToGroupID string
	//消息接收UserID
	ToUserID string
	//消息内容
	Content string
}

//get msgqueue with appid
func getMsgQueue(appId string) *MessageQueue {
	msgLock.RLock()
	defer msgLock.RUnlock()
	msgQueue, mok := msgQueueMap[appId]
	if mok {
		return msgQueue
	} else {
		return nil
	}
}

//push new message
//result:
//0:ok
//-10001:jsonunmar error
//-10002:this appid no have permission
func PushMessage(message string) (int, string) {
	//convert message json to LiveMessage
	var msg Message
	err_jsonunmar := json.Unmarshal([]byte(message), &msg)
	if err_jsonunmar != nil {
		logger.Error("PushMessage["+message+"]:jsonunmar error => ["+fmt.Sprintln(err_jsonunmar)+"]", LogTarget_Message)
		return -10001, "jsonunmar error"
	}
	msgQueue := getMsgQueue(msg.AppID)
	if msgQueue == nil {
		logger.Error("PushMessage["+message+"]:GetMsgQueue["+msg.AppID+"] Get => nil", LogTarget_Message)
		return -10002, "this appid no have permission"
	} else {
		msgQueue.MsgChan <- msg
		logger.Debug("PushMessage["+message+"] SetData => "+fmt.Sprintln(msg), LogTarget_Message)
		return 0, "ok"
	}
}

//put message from chan
func readMessage(appId string) (*Message, error) {
	logger.Debug("Begin ReadMessage["+appId+"] GetDate...", LogTarget_Message)
	var msg Message
	msgQueue := getMsgQueue(appId)
	if msgQueue == nil {
		logstr := "ReadMessage[" + appId + "]:GetMsgQueue[" + appId + "] Get => nil"
		logger.Error(logstr, LogTarget_Message)
		return nil, errors.New(logstr)
	} else {
		msg = <-msgQueue.MsgChan
		logger.Debug("ReadMessage["+appId+"] GetData => "+fmt.Sprintln(msg), LogTarget_Message)
		return &msg, nil
	}
}
