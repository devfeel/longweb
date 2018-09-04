package handlers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/longweb/message"
	"strconv"
	"strings"
)

func State(ctx dotweb.Context) error {
	appcount := len(message.AppPool)
	strhtml := "<html><body>{content}</body></html>"
	body := " AppCount => " + strconv.Itoa(appcount)
	body += "<br>  MaxClientIndex => " + strconv.FormatUint(message.GetMaxClientIndex(), 10)
	body += "<br>  ClientObjectCreateCount => " + strconv.FormatUint(message.GetTotalClientCreateCount(), 10)
	body += "<br>"
	app := ctx.QueryString("app")
	if app != "" {
		appinfo, exists := message.GetState_AppGroups(app)
		if !exists {
			body += "<br> not exists this app"
		} else {
			body += "<br>AppInfo [" + app + "] [UserGroupCount:" + strconv.Itoa(appinfo.GetGroupCount()) + "]"
			body += "<br>UserGroupList:"
			body += "<br>[" + strconv.Itoa(appinfo.GetState_TotalClientCount()) + "] =>"
			body += " [NW=" + strconv.Itoa(appinfo.GetState_WebSocketCount()-appinfo.GetState_AuthWebSocketCount()) + "]"
			body += " [AW=" + strconv.Itoa(appinfo.GetState_AuthWebSocketCount()) + "]"
			body += " [NL=" + strconv.Itoa(appinfo.GetState_LongPollCount()-appinfo.GetState_AuthLongPollCount()) + "]"
			body += " [AL=" + strconv.Itoa(appinfo.GetState_AuthLongPollCount()) + "]"
			body += "<br>"
			for _, v := range appinfo.GetState_UserGroups() {
				body += "<br>"
				line := "[" + v.GetGroupId() + " : " + strconv.Itoa(v.GetState_WebSocketClientCount()+v.GetState_LongPollClientCount()) + "] => "
				line += "[NW=" + strconv.Itoa(v.GetState_WebSocketClientCount()-v.GetState_AuthWebSocketClientCount()) + "] "
				line += "[AW=" + strconv.Itoa(v.GetState_AuthWebSocketClientCount()) + "] "
				line += "[NL=" + strconv.Itoa(v.GetState_LongPollClientCount()-v.GetState_AuthLongPollClientCount()) + "] "
				line += "[AL=" + strconv.Itoa(v.GetState_AuthLongPollClientCount()) + "]"

				body += line
			}
		}
	}
	strhtml = strings.Replace(strhtml, "{content}", body, 1)
	ctx.WriteHtml(strhtml)
	return nil
}

func StateData(ctx dotweb.Context) error {
	data := message.GetConnData()
	ctx.WriteJson(data)
	return nil
}
