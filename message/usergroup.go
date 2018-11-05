package message

import (
	"fmt"
	"github.com/devfeel/longweb/config"
	. "github.com/devfeel/longweb/const"
	"github.com/devfeel/longweb/framework/log"
	"github.com/influxdata/influxdb/pkg/slices"
	"strconv"
	"strings"
	"sync"
)

//表示用户组的集合
type AppGroups struct {
	appId  string
	groups map[string]*UserGroup
	mutex  *sync.RWMutex
}

//表示一组用户
type UserGroup struct {
	groupId          string
	websocketClients map[string]*UserClient
	longpollClients  map[string]*UserClient
	userMutex        *sync.RWMutex
}

//get usergroup with appid & groupid
func GetUserGroup(appId, groupId string) (*UserGroup, bool) {
	app, exists := GetAppGroups(appId)
	if !exists {
		return nil, false
	}

	app.mutex.RLock()
	group, mok := app.groups[groupId]
	app.mutex.RUnlock()
	return group, mok
}

//get appgroups with appid
func GetAppGroups(appId string) (*AppGroups, bool) {
	appLock.RLock()
	defer appLock.RUnlock()
	app, mok := AppPool[appId]
	return app, mok
}

//get appgroups with appid
func GetState_AppGroups(appId string) (*AppGroups, bool) {
	app, mok := AppPool[appId]
	return app, mok
}

//create a new client group
func NewUserGroup(groupId string) *UserGroup {
	return &UserGroup{groupId: groupId, userMutex: new(sync.RWMutex), websocketClients: make(map[string]*UserClient), longpollClients: make(map[string]*UserClient)}
}

//create a new groups
func NewAppGroups(appid string) *AppGroups {
	return &AppGroups{appId: appid, mutex: new(sync.RWMutex), groups: make(map[string]*UserGroup)}
}

//init appgroupinfo
func InitAppInfo() {
	appLock.Lock()
	defer appLock.Unlock()

	apps := config.GetAppList()
	for _, v := range apps {
		AppPool[v.AppID] = NewAppGroups(v.AppID)
	}
}

//get and init UserGroup with appid
func (app *AppGroups) GetAndInitUserGroup(appId, groupId string) *UserGroup {

	var group *UserGroup
	var mok bool

	app.mutex.RLock()
	group, mok = app.groups[groupId]
	app.mutex.RUnlock()

	//not exists, init
	if !mok {
		app.mutex.Lock()
		group, mok = app.groups[groupId]
		if !mok {
			group = NewUserGroup(groupId)
			app.groups[groupId] = group
		}
		app.mutex.Unlock()
	}

	return group
}

//get app's usergroup count
func (app *AppGroups) GetGroupCount() int {
	return len(app.groups)
}

//get app's total client count
func (app *AppGroups) GetState_TotalClientCount(groupIds ...string) int {
	total := 0
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	for _, g := range app.groups {
		if len(groupIds) > 0 && !slices.Exists(groupIds, g.groupId) {
			continue
		}
		total += g.GetState_WebSocketClientCount() + g.GetState_LongPollClientCount()
	}
	return total
}

//get app's websocke client count
func (app *AppGroups) GetState_WebSocketCount(groupIds ...string) int {
	total := 0
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	for _, g := range app.groups {
		if len(groupIds) > 0 && !slices.Exists(groupIds, g.groupId) {
			continue
		}
		total += g.GetState_WebSocketClientCount()
	}
	return total
}

//get app's auth websocke client count
func (app *AppGroups) GetState_AuthWebSocketCount(groupIds ...string) int {
	total := 0
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	for _, g := range app.groups {
		if len(groupIds) > 0 && !slices.Exists(groupIds, g.groupId) {
			continue
		}
		total += g.GetState_AuthWebSocketClientCount()
	}
	return total
}

//get app's longpoll client count
func (app *AppGroups) GetState_LongPollCount(groupIds ...string) int {
	total := 0
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	for _, g := range app.groups {
		if len(groupIds) > 0 && !slices.Exists(groupIds, g.groupId) {
			continue
		}
		total += g.GetState_LongPollClientCount()
	}
	return total
}

//get app's auth longpoll client count
func (app *AppGroups) GetState_AuthLongPollCount(groupIds ...string) int {
	total := 0
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	for _, g := range app.groups {
		if len(groupIds) > 0 && !slices.Exists(groupIds, g.groupId) {
			continue
		}
		total += g.GetState_AuthLongPollClientCount()
	}
	return total
}

//获取指定用户组
func (ag *AppGroups) GetUserGroup(groupId string) (*UserGroup, bool) {
	ag.mutex.RLock()
	defer ag.mutex.RUnlock()
	group, mok := ag.groups[groupId]
	return group, mok
}

//获取用户组列表
func (ag *AppGroups) GetState_UserGroups() map[string]*UserGroup {
	return ag.groups
}

//send a meeage for full app groups
//online groupid do nothing
//return send client count
func (ag *AppGroups) SendMessage(message *Message) int {
	clientcount := 0
	ag.mutex.RLock()
	defer ag.mutex.RUnlock()
	for _, group := range ag.groups {
		if strings.ToLower(group.groupId) != GroupID_Online {
			clientcount += group.SendMessage(message)
		}
	}
	return clientcount
}

//add new client into usergroup
func (ug *UserGroup) AddClient(client *UserClient) {
	ug.userMutex.Lock()
	if client.ConnType == ConnType_WebSocket {
		ug.websocketClients[client.UserID] = client
	} else if client.ConnType == ConnType_LongPoll {
		ug.longpollClients[client.UserID] = client
	}
	ug.userMutex.Unlock()
}

//delete a userclient
func (ug *UserGroup) DeleteClient(client *UserClient) {
	if client.ConnType == ConnType_WebSocket {
		ug.userMutex.Lock()
		delete(ug.websocketClients, client.UserID)
		ug.userMutex.Unlock()
	} else if client.ConnType == ConnType_LongPoll {
		ug.userMutex.Lock()
		delete(ug.longpollClients, client.UserID)
		ug.userMutex.Unlock()
	}
}

//send a meeage for full group
//return send client count
func (ug *UserGroup) SendMessage(message *Message) int {
	index := 0
	var needSend = true
	logger.Debug("UserClient:UserGroup->SendMessage("+fmt.Sprint(message)+") begin ["+strconv.Itoa(len(ug.websocketClients))+"]["+strconv.Itoa(len(ug.longpollClients))+"]", LogTarget_UserClient)
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(2)
	//发送websocket通道消息
	go func() {
		defer waitGroup.Done()
		ug.userMutex.RLock()
		for _, client := range ug.websocketClients {
			needSend = true
			if message.MessageLevel == MessageLevel_Auth {
				if !client.IsAuth {
					needSend = false
				}
			}
			if message.MessageLevel == MessageLevel_Normal {
				if client.IsAuth {
					needSend = false
				}
			}

			if needSend {
				index += 1
				client.SendMessage(message.Content)
			}
			if index%1000 == 0 {
				logger.Log("UserGroup["+ug.groupId+"]->SendMessage(WebSocket) Index => "+strconv.Itoa(index)+" success", LogTarget_UserClient, LogLevel_Debug)
			}
		}
		ug.userMutex.RUnlock()
	}()
	//发送longpool通道消息
	go func() {
		defer waitGroup.Done()
		ug.userMutex.RLock()
		for _, client := range ug.longpollClients {
			needSend = true
			if client.isHijackSend {
				needSend = false
			} else {
				if message.MessageLevel == MessageLevel_Auth {
					if !client.IsAuth {
						needSend = false
					}
				}
				if message.MessageLevel == MessageLevel_Normal {
					if client.IsAuth {
						needSend = false
					}
				}
			}

			if needSend {
				index += 1
				client.SendMessage(message.Content)
			}
			if index%1000 == 0 {
				logger.Log("UserGroup["+ug.groupId+"]->SendMessage(Hijack) Index => "+strconv.Itoa(index)+" success", LogTarget_UserClient, LogLevel_Debug)
			}
		}
		ug.userMutex.RUnlock()
	}()
	//等到通道发送完成
	waitGroup.Wait()
	logger.Debug("UserClient:UserGroup->SendMessage("+fmt.Sprint(message)+") end ["+strconv.Itoa(len(ug.websocketClients))+"]["+strconv.Itoa(len(ug.longpollClients))+"]", LogTarget_UserClient)
	return index
}

//获取指定UserId的用户客户端代理
func (ug *UserGroup) GetUserClient(userId string) (*UserClient, bool) {
	ug.userMutex.RLock()
	defer ug.userMutex.RUnlock()
	client, mok := ug.websocketClients[userId]
	if !mok {
		client, mok = ug.longpollClients[userId]
	}
	return client, mok
}

func (ug *UserGroup) GetGroupId() string {
	return ug.groupId
}

//get usergroup's websocketclient count
func (ug *UserGroup) GetState_WebSocketClientCount() int {
	return len(ug.websocketClients)
}

//get usergroup's auth websocketclient count
func (ug *UserGroup) GetState_AuthWebSocketClientCount() int {
	count := 0
	ug.userMutex.RLock()
	defer ug.userMutex.RUnlock()
	for _, client := range ug.websocketClients {
		if client.IsAuth {
			count += 1
		}
	}
	return count
}

//get usergroup's longpollclient count
func (ug *UserGroup) GetState_LongPollClientCount() int {
	return len(ug.longpollClients)
}

//get usergroup's auth longpollclient count
func (ug *UserGroup) GetState_AuthLongPollClientCount() int {
	count := 0
	ug.userMutex.RLock()
	defer ug.userMutex.RUnlock()
	for _, client := range ug.longpollClients {
		if client.IsAuth {
			count += 1
		}
	}
	return count
}
