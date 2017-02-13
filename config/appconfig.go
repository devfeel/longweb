package config

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"datapipe/websocket/framework/log"
)

var (
	CurrentConfig  AppConfig
	CurrentBaseDir string
	innerLogger    *logger.InnerLogger
	appMap         map[string]*AppInfo
	mutex          *sync.RWMutex
	allowIpMap     map[string]string
	allowIpMutex   *sync.RWMutex
)

func init() {
	//初始化读写锁
	mutex = new(sync.RWMutex)
	allowIpMutex = new(sync.RWMutex)
	innerLogger = logger.GetInnerLogger()
}

func SetBaseDir(baseDir string) {
	CurrentBaseDir = baseDir
}

//初始化配置文件
func InitConfig(configFile string) *AppConfig {
	innerLogger.Info("AppConfig::InitConfig 配置文件[" + configFile + "]开始...")
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		innerLogger.Warn("AppConfig::InitConfig 配置文件[" + configFile + "]无法解析 - " + err.Error())
		os.Exit(1)
	}

	var result AppConfig
	err = xml.Unmarshal(content, &result)
	if err != nil {
		innerLogger.Warn("AppConfig::InitConfig 配置文件[" + configFile + "]解析失败 - " + err.Error())
		os.Exit(1)
	}
	tmpMap := make(map[string]*AppInfo)
	for k, v := range result.Apps {
		tmpMap[v.AppID] = &result.Apps[k]
		b, _ := json.Marshal(&v)
		innerLogger.Info("AppConfig::InitConfig Load AppInfo => " + string(b))
	}

	//初始化App列表
	mutex.Lock()
	CurrentConfig = result
	appMap = tmpMap
	mutex.Unlock()

	//初始化AllowIP列表
	allowIpMutex.Lock()
	allowIpMap = make(map[string]string)
	innerLogger.Info("AppConfig::InitConfig => AllowIps => " + fmt.Sprint(result.AllowIps))
	for _, v := range result.AllowIps {
		innerLogger.Info("AppConfig::InitConfig => AddAllowIp -> [" + v + "]")
		allowIpMap[v] = v
	}
	allowIpMutex.Unlock()
	innerLogger.Info("AppConfig::InitConfig 配置文件[" + configFile + "]完成")

	return &CurrentConfig
}

func GetAppList() map[string]*AppInfo {
	return appMap
}

func GetAppInfo(appId string) (*AppInfo, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	app, exists := appMap[appId]
	return app, exists
}

//检测IP是否被允许访问
func CheckAllowIP(ip string) bool {
	allowIpMutex.RLock()
	_, exists := allowIpMap[ip]
	allowIpMutex.RUnlock()
	//logger.Log("CheckAllowIP["+ip+"] => "+strconv.FormatBool(exists), LogTarget_Default, LogLevel_Debug)
	return exists
}
