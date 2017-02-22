package message

import (
	"github.com/devfeel/longweb/config"
	"github.com/devfeel/longweb/framework/http"
	"github.com/devfeel/longweb/framework/json"
	"strconv"
)

/*检查应用鉴权结果
Author: panxinming
CreateTime: 2017-02-07 13:00
入参：
app：应用配置信息
appId：应用编号
groupid：用户组编号，应用自定义
userid：用户编号，应用需保证userid在同Appid下的唯一性
token：鉴权token
返回值：
0 成功
-101001：no permission connect! => appid[" + appId + "] no exists
-101002：no permission connect! => check token has an error => [http-request-error]
-101003：no permission connect! => check token has an error => [json-parse-error]
-101004：no permission connect! => check token result =>
-101005：no permission connect! => check token: appid|groupid|userid not match

update log:
1、初始版本 --2017-02-07 13:00 by pxm
*/
func CheckAuthToken(app *config.AppInfo, appId, groupId, userId, token string) (RetCode int, RetMsg string) {
	RetCode = 0
	RetMsg = ""
	//check AuthApi
	if app.AuthApi == "" {
		RetCode = -101001
		RetMsg = "no permission connect! => appid[" + appId + "] no exists"
		return
	}

	//check token from app
	type AuthResponse struct {
		RetCode int
		RetMsg  string
		AppID   string
		GroupID string
		UserID  string
	}

	targetUrl := app.AuthApi + "?appid=" + appId + "&groupid=" + groupId + "&userid=" + userId + "&token=" + token
	authbody, _, _, autherr := httputil.HttpGet(targetUrl)
	if autherr != nil {
		RetCode = -101002
		RetMsg = "no permission connect! => check token has an error => [http-request-error] " + autherr.Error()
		return
	} else {
		var authRes AuthResponse
		jsonerr := jsonutil.Unmarshal(authbody, &authRes)
		if jsonerr != nil {
			RetCode = -101003
			RetMsg = "no permission connect! => check token has an error => [json-parse-error] " + jsonerr.Error()
			return
		} else {
			if authRes.RetCode != 0 {
				RetCode = -101004
				RetMsg = "no permission connect! => check token result => " + strconv.Itoa(authRes.RetCode)
				return
			} else {
				if authRes.AppID != appId || authRes.GroupID != groupId || authRes.UserID != userId {
					RetCode = -101005
					RetMsg = "no permission connect! => check token: appid|groupid|userid not match"
					return
				}
			}
		}
	}
	return
}
