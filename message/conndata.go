package message

import "time"

type AppConnData struct {
	AppID           string
	TotalCount      int
	NormalWebsocket int
	AuthWebsocket   int
	NormalLongPoll  int
	AuthLongPoll    int
}
type ConnData struct {
	Apps     []AppConnData
	DataTime time.Time
}

const OnlineGroupID = "online"

//获取连接信息
func GetConnData() *ConnData {
	data := &ConnData{}
	data.DataTime = time.Now()
	data.Apps = make([]AppConnData, 0)
	totalData := AppConnData{AppID: "total"}
	for appid, appinfo := range AppPool {
		tmpData := AppConnData{
			AppID:           appid,
			TotalCount:      appinfo.GetTotalClientCount(OnlineGroupID),
			NormalWebsocket: appinfo.GetWebSocketCount(OnlineGroupID) - appinfo.GetAuthWebSocketCount(OnlineGroupID),
			AuthWebsocket:   appinfo.GetAuthWebSocketCount(OnlineGroupID),
			NormalLongPoll:  appinfo.GetLongPollCount(OnlineGroupID) - appinfo.GetAuthLongPollCount(OnlineGroupID),
			AuthLongPoll:    appinfo.GetAuthLongPollCount(OnlineGroupID),
		}
		data.Apps = append(data.Apps, tmpData)
		totalData.TotalCount += tmpData.TotalCount
		totalData.NormalWebsocket += tmpData.NormalWebsocket
		totalData.AuthWebsocket += tmpData.AuthWebsocket
		totalData.NormalLongPoll += tmpData.NormalLongPoll
		totalData.AuthLongPoll += tmpData.AuthLongPoll
	}
	data.Apps = append(data.Apps, totalData)
	return data
}
