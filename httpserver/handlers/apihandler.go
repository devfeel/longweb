package handlers

import (
	"devfeel/dotweb"
	"devfeel/longweb/config"
	. "devfeel/longweb/message"
)

//推送消息
//result:
//0:ok
//-1:not allow ip
//-2:message is empty
//-10001:message format error
//-10002:this appid no have permission
func SendMessage(ctx *dotweb.HttpContext) {
	type retJson struct {
		RetCode int
		RetMsg  string
	}

	result := &retJson{RetCode: 0, RetMsg: ""}

	//check allow ip
	remoteIp := ctx.RemoteIP()
	if !config.CheckAllowIP(remoteIp) {
		result.RetCode = -1
		result.RetMsg = "not allow ip"
	}

	//push message
	if result.RetCode == 0 {
		message := string(ctx.PostBody())
		if message != "" {
			result.RetCode, result.RetMsg = PushMessage(message)
		} else {
			result.RetCode = -2
			result.RetMsg = "message is empty"
		}
	}
	ctx.WriteJson(result)
}
