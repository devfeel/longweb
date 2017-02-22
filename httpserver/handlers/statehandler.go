package handlers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/longweb/message"
	"strconv"
	"strings"
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
