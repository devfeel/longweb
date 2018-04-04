package handlers

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/longweb/const"
	"github.com/devfeel/longweb/framework/file"
	"github.com/devfeel/longweb/framework/json"
	"html/template"
	"runtime"
)

func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")

	//ctx.WriteString("welcome to websocket proxy<br>" + ctx.Request.RequestURI + "<br>" + fmt.Sprintln(ctx.Request))
	ctx.WriteString("welcome to websocket proxy | version=" + constdefine.Version)
	return nil
}

func Memstate(ctx dotweb.Context) error {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)
	ctx.WriteString(fmt.Sprint(stats))
	return nil
}

func Test(ctx dotweb.Context) error {
	filePath := fileutil.GetCurrentDirectory()
	filePath = filePath + "/test.html"
	tmpl, err := template.New("test.html").ParseFiles(filePath)
	if err != nil {
		ctx.WriteString("version template Parse error => " + err.Error())
		return err
	}
	viewdata := make(map[string]string)
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf8")
	err = tmpl.Execute(ctx.Response().Writer(), viewdata)
	if err != nil {
		ctx.WriteString("version template Execute error => " + err.Error())
		return err
	}
	return nil
}

func TestAuth(ctx dotweb.Context) error {
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
	return ctx.WriteString(body)
}

func TestMessage(ctx dotweb.Context) error {
	return ctx.WriteString("")
}
