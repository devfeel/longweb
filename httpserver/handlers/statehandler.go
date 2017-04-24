package handlers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/longweb/message"
	"strconv"
	"strings"
	"time"
)

func State(ctx *dotweb.HttpContext) {
	appcount := len(message.AppPool)
	strhtml := "<html><body>{content}</body></html>"
	body := " AppCount => " + strconv.Itoa(appcount)
	body += "<br>  MaxClientIndex => " + strconv.FormatUint(message.GetMaxClientIndex(), 10)
	body += "<br>  ClientObjectCreateCount => " + strconv.FormatUint(message.GetTotalClientCreateCount(), 10)
	body += "<br>"
	app := ctx.QueryString("app")
	if app != "" {
		appinfo, exists := message.GetAppGroups(app)
		if !exists {
			body += "<br> not exists this app"
		} else {
			body += "<br>AppInfo [" + app + "] [UserGroupCount:" + strconv.Itoa(appinfo.GetGroupCount()) + "]"
			body += "<br>UserGroupList:"
			body += "<br>[" + strconv.Itoa(appinfo.GetTotalClientCount()) + "] =>"
			body += " [NW=" + strconv.Itoa(appinfo.GetWebSocketCount()-appinfo.GetAuthWebSocketCount()) + "]"
			body += " [AW=" + strconv.Itoa(appinfo.GetAuthWebSocketCount()) + "]"
			body += " [NL=" + strconv.Itoa(appinfo.GetLongPollCount()-appinfo.GetAuthLongPollCount()) + "]"
			body += " [AL=" + strconv.Itoa(appinfo.GetAuthLongPollCount()) + "]"
			body += "<br>"
			for _, v := range appinfo.GetUserGroups() {
				body += "<br>"
				line := "[" + v.GetGroupId() + " : " + strconv.Itoa(v.GetWebSocketClientCount()+v.GetLongPollClientCount()) + "] => "
				line += "[NW=" + strconv.Itoa(v.GetWebSocketClientCount()-v.GetAuthWebSocketClientCount()) + "] "
				line += "[AW=" + strconv.Itoa(v.GetAuthWebSocketClientCount()) + "] "
				line += "[NL=" + strconv.Itoa(v.GetLongPollClientCount()-v.GetAuthLongPollClientCount()) + "] "
				line += "[AL=" + strconv.Itoa(v.GetAuthLongPollClientCount()) + "]"

				body += line
			}
		}
	}
	strhtml = strings.Replace(strhtml, "{content}", body, 1)
	ctx.WriteString(strhtml)
}

func StateData(ctx *dotweb.HttpContext) {
	type AppData struct {
		AppID           string
		TotalCount      int
		NormalWebsocket int
		AuthWebsocket   int
		NormalLongPoll  int
		AuthLongPoll    int
	}
	type ConnData struct {
		Apps     []AppData
		DataTime time.Time
	}

	data := ConnData{}
	data.DataTime = time.Now()
	data.Apps = make([]AppData, 0)
	totalData := AppData{AppID: "total"}
	for appid, appinfo := range message.AppPool {
		tmpData := AppData{
			AppID:           appid,
			TotalCount:      appinfo.GetTotalClientCount(),
			NormalWebsocket: appinfo.GetWebSocketCount() - appinfo.GetAuthWebSocketCount(),
			AuthWebsocket:   appinfo.GetAuthWebSocketCount(),
			NormalLongPoll:  appinfo.GetLongPollCount() - appinfo.GetAuthLongPollCount(),
			AuthLongPoll:    appinfo.GetAuthLongPollCount(),
		}
		data.Apps = append(data.Apps, tmpData)
		totalData.TotalCount += tmpData.TotalCount
		totalData.NormalWebsocket += tmpData.NormalWebsocket
		totalData.AuthWebsocket += tmpData.AuthWebsocket
		totalData.NormalLongPoll += tmpData.NormalLongPoll
		totalData.AuthLongPoll += tmpData.AuthLongPoll
	}
	data.Apps = append(data.Apps, totalData)
	ctx.WriteJson(data)
}
