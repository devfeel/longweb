package tasks

import (
	"errors"
	"github.com/devfeel/longweb/config"
	"github.com/devfeel/longweb/exception"
	"github.com/devfeel/longweb/framework/log"
	"github.com/devfeel/longweb/message"
	"github.com/devfeel/longweb/repository"
	"github.com/devfeel/dottask"
	"time"
)

func Task_SyncOnlineData(context *task.TaskContext) error {
	logTitle := "Tasks:Task_MonitorChartData:"

	defer func() {
		if err := recover(); err != nil {
			ex := exception.CatchError(logTitle+"defer-recover:", err)
			logger.Error(ex.GetDefaultLogString(), context.TaskID)
		}
	}()

	//从本地获取数据
	data := message.GetConnData()
	influxdbImpl := new(repository.InfluxdbImpl)
	influxConf := config.CurrentConfig.SyncNode.InfluxdbInfo
	if influxConf == nil || influxConf.ServerIP == "" {
		return errors.New("未正确配置InfluxDB")
	}
	influxdbImpl.SetConn(influxConf.ServerIP, influxConf.DBName, influxConf.UserName, influxConf.Password)
	//拼接influxdb结构
	influxData := repository.NewInfluxdbData()
	influxData.TableName = "appstat"
	for _, v := range data.Apps {
		influxData.Tags["id"] = v.AppID
		influxData.Tags["name"] = ""
		influxData.Tags["url"] = ""
		influxData.Fields["total"] = v.TotalCount
		influxData.Fields["normalwebsocket"] = v.NormalWebsocket
		influxData.Fields["authwebsocket"] = v.AuthWebsocket
		influxData.Fields["normallongpoll"] = v.NormalLongPoll
		influxData.Fields["authlongpoll"] = v.AuthLongPoll
		influxData.Time = time.Now()
		influxdbImpl.InsertData(influxData)
	}
	return nil
}
