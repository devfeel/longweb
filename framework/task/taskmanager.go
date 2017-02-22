package task

import (
	"github.com/devfeel/longweb/framework/log"
	"sync"
	"time"
)

var (
	taskMap     map[string]*(TaskInfo)
	innerLogger *logger.InnerLogger
	mutex       *sync.RWMutex
)

//task info
type TaskInfo struct {
	Name string
	//创建Task时，可通过该属性传递需要定制的信息
	TaskData  interface{}
	handler   TaskHandler
	timeTiker *(time.Ticker)
}

type TaskHandler func(*TaskInfo)

func init() {
	mutex = new(sync.RWMutex)
	taskMap = make(map[string]*(TaskInfo))
	innerLogger = logger.GetInnerLogger()
}

//停止指定Task执行
func (t *TaskInfo) Stop() {
	t.timeTiker.Stop()
}

//启动指定Task执行
func (t *TaskInfo) Start() {
	t.timeTiker = time.NewTicker(1 * time.Millisecond)
	for {
		select {
		case <-t.timeTiker.C:
			t.handler(t)
		}
	}
}

//传入Name创建Taskinfo
func NewTask(name string, handler TaskHandler) *TaskInfo {
	t := &TaskInfo{
		Name:    name,
		handler: handler,
	}
	return t
}

//结束所有Task
func StopAllTask() {
	innerLogger.Info("Task:StopAllTask begin...")
	for _, v := range taskMap {
		innerLogger.Info("Task:StopAllTask:StopTask => " + v.Name)
		v.Stop()
	}
	innerLogger.Info("Task::StopAllTask end[" + string(len(taskMap)) + "]")
}

//开启所有Task执行
func StartAllTask() {
	innerLogger.Info("Task:StartAllTask begin...")
	for _, v := range taskMap {
		innerLogger.Info("Task:StartAllTask:StartTask => " + v.Name)
		go v.Start()
	}
	innerLogger.Info("Task:StartAllTask end")
}

func RegisterTask(task *TaskInfo) {
	mutex.Lock()
	defer mutex.Unlock()
	taskMap[task.Name] = task
}

func RemoveTask(taskName string) {
	mutex.Lock()
	defer mutex.Unlock()
	t, exists := taskMap[taskName]
	if exists {
		t.Stop()
		delete(taskMap, taskName)
	}
}

func ReStartAllTask() {
	innerLogger.Info("Task::ReStartAllTask begin...")
	StopAllTask()
	StartAllTask()
	innerLogger.Info("Task::ReStartAllTask end")
}
