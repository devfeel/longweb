package logger

import (
	"os"
	"strings"
	"syscall"
	"time"
)

var innerLogRootPath string

type InnerLogger struct {
}

var singleInnerLogger *InnerLogger

func GetInnerLogger() *InnerLogger {
	if singleInnerLogger == nil {
		singleInnerLogger = new(InnerLogger)
	}
	return singleInnerLogger
}

func (logger *InnerLogger) Debug(log string) {
	logger.innerWriteLog(log, "debug")
}

func (logger *InnerLogger) Info(log string) {
	logger.innerWriteLog(log, "info")
}

func (logger *InnerLogger) Warn(log string) {
	logger.innerWriteLog(log, "warn")
}

func (logger *InnerLogger) Error(log string) {
	logger.innerWriteLog(log, "error")
}

//开启日志处理器
func StartInnerLogHandler(rootPath string) {
	//设置日志根目录
	innerLogRootPath = rootPath
	if !strings.HasSuffix(innerLogRootPath, "/") {
		innerLogRootPath = innerLogRootPath + "/"
	}
}

func (logger *InnerLogger) innerWriteLog(log string, level string) {
	filePath := innerLogRootPath + "/innerlogs/"
	switch level {
	case "debug":
		filePath = filePath + "innerlog_debug_" + time.Now().Format(defaultDateFormatForFileName) + ".log"
	case "info":
		filePath = filePath + "innerlog_info_" + time.Now().Format(defaultDateFormatForFileName) + ".log"
	case "warn":
		filePath = filePath + "innerlog_warn_" + time.Now().Format(defaultDateFormatForFileName) + ".log"
	case "error":
		filePath = filePath + "innerlog_error_" + time.Now().Format(defaultDateFormatForFileName) + ".log"
		break
	}
	log = time.Now().Format(defaultFullTimeLayout) + " " + log
	logger.innerWriteFile(filePath, log)
}

func (logger *InnerLogger) innerWriteFile(logFile string, log string) {
	var mode os.FileMode
	flag := syscall.O_RDWR | syscall.O_APPEND | syscall.O_CREAT
	mode = 0666
	logstr := log + "\r\n"
	file, err := os.OpenFile(logFile, flag, mode)
	defer file.Close()
	if err != nil {
		return
	}
	file.WriteString(logstr)
}
