package handlers

import (
	"devfeel/dotweb"
	"devfeel/longweb/message"
	"strconv"
	"strings"
)

func State(ctx *dotweb.HttpContext) {
	appcount := len(message.AppPool)
	strhtml := "<html><body>{content}</body></html>"
	body := " AppCount => " + strconv.Itoa(appcount)
	body += "<br>  MaxClientIndex => " + strconv.FormatUint(message.GetMaxClientIndex(), 10)
	body += "<br>  ClientObjectCreateCount => " + strconv.FormatUint(message.GetTotalClientCreateCount(), 10)
	app := ctx.QueryString("app")
	if app != "" {
		appinfo, exists := message.GetAppGroups(app)
		if !exists {
			body += "<br> not exists this app"
		} else {
			body += "<br> UserGroupCount => " + strconv.Itoa(appinfo.GetGroupCount())
			body += "<br> UserClientCount => " + strconv.Itoa(appinfo.GetTotalClientCount())
			body += "<br><br>[" + app + "] => UserGroupList:"
			body += "<br>[GroupID:TotalCount] => [NW=\"NormalWebSocket\"] [AW=\"AuthWebSocket\"] [NL=\"NormalLongPoll\"] [AL=\"AuthLongPoll\"]"
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
