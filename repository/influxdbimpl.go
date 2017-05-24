package repository

import (
	"fmt"
	"github.com/devfeel/longweb/framework/log"
	"github.com/influxdata/influxdb/client/v2"
	"time"
)

type InfluxdbImpl struct {
	serverIp  string
	username  string
	password  string
	defaultDB string
}

type InfluxdbData struct {
	TableName string
	Tags      map[string]string
	Fields    map[string]interface{}
	Time      time.Time
}

func NewInfluxdbData() *InfluxdbData {
	return &InfluxdbData{
		Tags:   make(map[string]string),
		Fields: make(map[string]interface{}),
	}
}

const logTarget = "InfluxdbImpl"

//设置连接属性
func (da *InfluxdbImpl) SetConn(serverIp, dbName, username, password string) {
	da.serverIp = serverIp
	da.defaultDB = dbName
	da.username = username
	da.password = password
}

//insert data into influxdb
func (da *InfluxdbImpl) InsertData(data *InfluxdbData) error {
	logtitle := "repository.InfluxdbImpl:InsertData:"
	/*c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     da.serverIp,
		Username: da.username,
		Password: da.password,
	})*/

	c, err := client.NewUDPClient(client.UDPConfig{
		Addr: da.serverIp,
	})

	if err != nil {
		logger.Error(logtitle+"NewHTTPClient["+da.serverIp+"] error - "+err.Error(), logTarget)
		logger.Debug(logtitle+"NewHTTPClient["+da.serverIp+"] error", logTarget)
		return err
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  da.defaultDB,
		Precision: "s",
	})
	if err != nil {
		logger.Error(logtitle+"NewBatchPoints["+da.defaultDB+"] error - "+err.Error(), logTarget)
		logger.Debug(logtitle+"NewBatchPoints["+da.defaultDB+"] error", logTarget)
		return err
	}

	pt, err := client.NewPoint(data.TableName, data.Tags, data.Fields, data.Time)
	if err != nil {
		logger.Error(logtitle+"NewPoint["+fmt.Sprint(data)+"] error - "+err.Error(), logTarget)
		logger.Debug(logtitle+"NewPoint["+fmt.Sprint(data)+"] error", logTarget)
		return err
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		logger.Error(logtitle+"Write["+fmt.Sprint(data)+"] error - "+err.Error(), logTarget)
		logger.Debug(logtitle+"Write["+fmt.Sprint(data)+"] error", logTarget)
		return err
	} else {
		logger.Debug(logtitle+"Write["+fmt.Sprint(data)+"] success", logTarget)
	}
	return nil
}
