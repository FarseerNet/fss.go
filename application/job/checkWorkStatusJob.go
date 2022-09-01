package job

import (
	"context"
	"fmt"
	"fss/domain/clients/client"
	"fss/domain/log"
	"fss/domain/tasks/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/tasks"
	"time"
)

// RegisterCheckWorkStatusJob 检测进行中状态的任务
func RegisterCheckWorkStatusJob() {
	tasks.Run("FSS.CheckWorkStatus", 30*time.Second, checkWorkStatusJob, context.Background())
}

func checkWorkStatusJob(context *tasks.TaskContext) {
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	for _, taskGroupDO := range taskGroupRepository.ToSchedulerWorkingList().ToArray() {
		checkTaskGroup(taskGroupDO, taskGroupRepository)
	}
}

func checkTaskGroup(taskGroup taskGroup.DomainObject, taskGroupRepository taskGroup.Repository) {
	taskGroup = taskGroupRepository.ToEntity(taskGroup.Id)
	if taskGroup.Task.Id < 1 {
		taskGroup.CreateTask()
		taskGroupRepository.Save(taskGroup)
		return
	}

	// 任务不存在
	if taskGroup.Task.Client.Id > 0 {
		clientDO := container.Resolve[client.Repository]().ToEntity(taskGroup.Task.Client.Id)
		if clientDO.Id < 1 {
			message := fmt.Sprint("【客户端不存在】", taskGroup.Task.Client.Id, "，强制下线客户端")
			log.TaskLogAddService(taskGroup.Id, taskGroup.JobName, taskGroup.Caption, eumLogLevel.Warning, message)
			taskGroup.Cancel()
			taskGroupRepository.Save(taskGroup)
			return
		}
	}

	taskGroup.CheckClientOffline()
}
