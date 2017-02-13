package handlers

import (
	"devfeel/dotweb"
	"devfeel/longweb/const"
	"devfeel/longweb/framework/file"
	"devfeel/longweb/framework/json"
	"fmt"
	"html/template"
	"runtime"
)

func Index(ctx *dotweb.HttpContext) {
	ctx.Response.Header().Set("Content-Type", "text/html; charset=utf-8")

	//ctx.WriteString("welcome to websocket proxy<br>" + ctx.Request.RequestURI + "<br>" + fmt.Sprintln(ctx.Request))
	ctx.WriteString("welcome to websocket proxy | version=" + constdefine.Version)
}

func Memstate(ctx *dotweb.HttpContext) {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)
	ctx.WriteString(fmt.Sprint(stats))
}

func Test(ctx *dotweb.HttpContext) {
	filePath := fileutil.GetCurrentDirectory()
	filePath = filePath + "/test.html"
	tmpl, err := template.New("test.html").ParseFiles(filePath)
	if err != nil {
		ctx.WriteString("version template Parse error => " + err.Error())
		return
	}
	viewdata := make(map[string]string)
	ctx.Response.Header().Set("Content-Type", "text/html; charset=utf8")
	err = tmpl.Execute(ctx.Response.Writer(), viewdata)
	if err != nil {
		ctx.WriteString("version template Execute error => " + err.Error())
		return
	}
}

func TestAuth(ctx *dotweb.HttpContext) {
	appId := ctx.QueryString("appid")
	groupId := ctx.QueryString("groupid")
	userId := ctx.QueryString("userid")
	token := ctx.QueryString("token")

	type AuthResponse struct {
		RetCode int
		RetMsg  string
		AppID   string
		GroupID string
		UserID  string
	}

	var res AuthResponse
	res.RetCode = 0
	res.RetMsg = token
	res.AppID = appId
	res.GroupID = groupId
	res.UserID = userId

	body := jsonutil.GetJsonString(res)
	ctx.WriteString(body)
}

func TestMessage(ctx *dotweb.HttpContext) {
	ctx.WriteString("")
}
