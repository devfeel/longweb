package task

import (
	"github.com/devfeel/longweb/task/tasks"
	"github.com/devfeel/task"
)

var service *task.TaskService

func RegisterTaskHandler(service *task.TaskService) {
	service.RegisterHandler("synconlinedata", tasks.Task_SyncOnlineData)
}

func StartTaskService(configFile string) {
	//step 1: init new task service
	service = task.StartNewService()

	//step 2: register all task handler
	RegisterTaskHandler(service)

	//step 3: load config file
	service.LoadConfig(configFile)

	//step 4: start all task
	service.StartAllTask()
}
