package handlers

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/longweb/config"
	. "github.com/devfeel/longweb/const"
	"github.com/devfeel/longweb/framework/http"
	"github.com/devfeel/longweb/framework/log"
	. "github.com/devfeel/longweb/message"
	"net/url"
	"strconv"
)

/*longpoll统一处理入口 - 兼容Hijack与HttpRequest
Author: panxinming
CreateTime: 2017-01-25 13:00
入参：
jsonpcallback: 应用定义jsonp-callback函数名，如果不传入默认为callback
appid：应用编号，统一申请
groupid：用户组编号，应用自定义
userid：用户编号，应用需保证userid在同Appid下的唯一性
querykey：透传key，会透传给应用messageapi，一般用于决定是否有需要马上返回的数据
返回值：
0 成功
-100001：appid、groupid、querykey不能为空
-100002：指定appid不存在
-100003：注册失败
-100009：超时或其他异常，一般需重新发起请求
-200009：应用返回异常

//鉴权相关返回值
-101001：no permission connect! => appid[" + appId + "] no exists
-101002：no permission connect! => check token has an error => [http-request-error]
-101003：no permission connect! => check token has an error => [json-parse-error]
-101004：no permission connect! => check token result =>
-101005：no permission connect! => check token: appid|groupid|userid not match

update log:
1、由于无法支持跨域，调整Hijack为HttpRequest --2017-02-04 by pxm
2、完善代码，Hijack模式支持跨域 --2017-02-05 by pxm
3、整合代码，兼容Hijack与HttpRequest --2017-02-07 by pxm
4、增加jsonp支持 --2017-02-08 17:00
*/
func OnPolling(ctx dotweb.Context) error {

	type ResponseJson struct {
		RetCode int
		RetMsg  string
		Message string
	}
	var resJson ResponseJson
	resJson.RetCode = 0
	resJson.RetMsg = "ok"
	resJson.Message = ""

	jsonpcallback := ctx.QueryString("jsonpcallback")
	appId := ctx.QueryString("appid")
	groupId := ctx.QueryString("groupid")
	userId := ctx.QueryString("userid")
	from := ctx.QueryString("from")
	querykey := ctx.QueryString("querykey")
	token := ctx.QueryString("token")
	//针对jsonp未传情况，默认为callback
	if jsonpcallback == "" {
		jsonpcallback = "callback"
	}

	defer func() {
		if ctx.IsHijack() {
			ctx.HijackConn().Close()
		}

	}()

	logTitle := "[OnPolling][" + ctx.Request().Url() + "]"
	if ctx.IsHijack() {
		logTitle += "[HiJack]"

	}
	logTitle += " "

	//处理跨域支持
	ctx.Response().SetHeader("Access-Control-Allow-Origin", "*")
	ctx.Response().SetHeader("P3P", "CP=\"CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR\"")

	logger.Debug(logTitle+"["+ctx.Request().Url()+"] connect [RemoteIp:"+ctx.RemoteIP()+"]", LogTarget_HttpRequest)

	if appId == "" || groupId == "" || querykey == "" {
		resJson.RetCode = -100001
		resJson.RetMsg = "not supported querystring => " + ctx.Request().RequestURI
		logger.Warn(logTitle+resJson.RetMsg, LogTarget_LongPoll)
		ctx.WriteJsonp(jsonpcallback, resJson)
		return nil
	}

	app, exists := config.GetAppInfo(appId)
	if !exists {
		resJson.RetCode = -100002
		resJson.RetMsg = "no permission connect! => appid[" + appId + "] no exists"
		logger.Warn(logTitle+resJson.RetMsg, LogTarget_LongPoll)
		ctx.WriteJsonp(jsonpcallback, resJson)
		return nil
	}

	//如果需要验证token
	if token != "" {
		retCode, retMsg := CheckAuthToken(app, appId, groupId, userId, token)
		if retCode != 0 {
			resJson.RetCode = retCode
			resJson.RetMsg = retMsg
			logger.Warn(logTitle+resJson.RetMsg, LogTarget_LongPoll)
			ctx.WriteJsonp(jsonpcallback, resJson)
			return nil
		}
	}

	isAuth := false
	if token != "" {
		isAuth = true
	}

	client := NewClient(appId, userId, groupId, from, isAuth, nil, ctx)
	client.TimeOut = app.TimeOut
	defer RemoveClient(client)

	//注册客户端
	_, regCode := RegisterClient(client)
	if regCode != 0 {
		resJson.RetCode = -100003
		resJson.RetMsg = "no permission connect! => " + strconv.Itoa(regCode)
		logger.Warn(logTitle+"["+client.GetClientInfo()+"] "+resJson.RetMsg, LogTarget_LongPoll)
		ctx.WriteJsonp(jsonpcallback, resJson)
		return nil
	}

	//如果MessageApi未配置，则忽略首次查询
	if app.MessageApi != "" {
		//get now data from app
		targetQuery := ""
		sourceQuery := "appid=" + appId + "&groupid=" + groupId + "&userid=" + userId + "&querykey=" + querykey
		parseQuery, errParse := url.ParseQuery(sourceQuery)
		if errParse != nil {
			targetQuery = sourceQuery
			logger.Warn(logTitle+"["+client.GetClientInfo()+"] UrlParse["+sourceQuery+"] error => "+errParse.Error(), LogTarget_LongPoll)
		} else {
			targetQuery = parseQuery.Encode()
		}
		targetUrl := app.MessageApi + "?" + targetQuery
		body, _, _, err := httputil.HttpGet(targetUrl)
		if err != nil {
			resJson.RetCode = -200009
			resJson.RetMsg = body
			resJson.Message = err.Error()
			logger.Warn(logTitle+"["+client.GetClientInfo()+"] HttpGet["+targetUrl+"] error => "+err.Error(), LogTarget_LongPoll)
			ctx.WriteJsonp(jsonpcallback, resJson)
			return nil
		} else {
			logger.Debug(logTitle+"["+client.GetClientInfo()+"] HttpGet["+targetUrl+"] success return => "+body, LogTarget_LongPoll)
			if body != "" {
				resJson.RetCode = 0
				resJson.RetMsg = "ok"
				resJson.Message = body
				ctx.WriteJsonp(jsonpcallback, resJson)
				return nil
			}
		}
	}

	//等待固定时间，等待消息推送
	if strMsg, err := client.ReadMessage(); err != nil {
		resJson.RetCode = -100009
		resJson.RetMsg = err.Error()
		logger.Warn(logTitle+"["+client.GetClientInfo()+"] ReadMessage error => "+err.Error(), LogTarget_LongPoll)
		ctx.WriteJsonp(jsonpcallback, resJson)
		return nil
	} else {
		resJson.RetCode = 0
		resJson.RetMsg = "ok"
		resJson.Message = strMsg
		logger.Debug(logTitle+"["+client.GetClientInfo()+"] readmessage return => "+strMsg, LogTarget_LongPoll)
		ctx.WriteJsonp(jsonpcallback, resJson)
		return nil
	}
}
