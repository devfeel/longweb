package config

import (
	"encoding/xml"
)

//配置信息
type AppConfig struct {
	XMLName    xml.Name   `xml:"config"`
	Log        Log        `xml:"log"`
	Apps       []AppInfo  `xml:"apps>app"`
	HttpServer HttpServer `xml:"httpserver"`
	AllowIps   []string   `xml:"allowips>ip"`
}

//全局配置
type HttpServer struct {
	HttpPort  int `xml:"httpport,attr"`
	PProfPort int `xml:"pprofport,attr"`
}

//log配置
type Log struct {
	FilePath string `xml:"filepath,attr"`
}

//app配置
type AppInfo struct {
	AppID   string `xml:"appid,attr"`
	AppName string `xml:"appname,attr"`
	//域名
	Domain     string `xml:"domain,attr"`
	MessageApi string `xml:"messageapi,attr"`
	//鉴权Api
	AuthApi string `xml:"authapi,attr"`
	TimeOut int64  `xml:"timeout,attr"`
}
