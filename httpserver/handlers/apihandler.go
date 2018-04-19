package handlers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/longweb/config"
	. "github.com/devfeel/longweb/message"
)

//推送消息
//result:
//0:ok
//-1:not allow ip
//-2:message is empty
//-10001:message format error
//-10002:this appid no have permission
func SendMessage(ctx dotweb.Context) error {
	type retJson struct {
		RetCode int
		RetMsg  string
	}

	result := &retJson{RetCode: 0, RetMsg: ""}

	//check allow ip
	remoteIp := ctx.RemoteIP()
	if !config.CheckAllowIP(remoteIp) {
		result.RetCode = -1
		result.RetMsg = "not allow ip =>" + remoteIp
	}

	//push message
	if result.RetCode == 0 {
		message := string(ctx.Request().PostBody())
		if message != "" {
			result.RetCode, result.RetMsg = PushMessage(message)
		} else {
			result.RetCode = -2
			result.RetMsg = "message is empty"
		}
	}
	return ctx.WriteJson(result)
}
