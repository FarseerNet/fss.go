package job

import (
	"context"
	"fmt"
	"fss/domain/clients/client"
	"fss/domain/log"
	"fss/domain/tasks/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/exception"
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

func checkTaskGroup(taskGroupDO taskGroup.DomainObject, taskGroupRepository taskGroup.Repository) {
	taskGroupDO = taskGroupRepository.ToEntity(taskGroupDO.Id)
	if taskGroupDO.IsNull() {
		return
	}

	if taskGroupDO.Task.IsNull() || taskGroupDO.Task.IsFinish() {
		taskGroupDO.CreateTask()
		taskGroupRepository.Save(taskGroupDO)
		return
	}

	// 任务不存在
	if taskGroupDO.Task.Client.Id > 0 {
		clientDO := container.Resolve[client.Repository]().ToEntity(taskGroupDO.Task.Client.Id)
		if clientDO.Id < 1 {
			message := fmt.Sprint("【客户端不存在】", taskGroupDO.Task.Client.Id, "，强制下线客户端")
			log.TaskLogAddService(taskGroupDO.Id, taskGroupDO.JobName, taskGroupDO.Caption, eumLogLevel.Warning, message)
			taskGroupDO.Cancel()
			taskGroupRepository.Save(taskGroupDO)
			return
		}
	}

	exception.Try(func() {
		taskGroupDO.CheckClientOffline()
	}).CatchRefuseException(func(exp *exception.RefuseException) {
		log.TaskLogAddService(taskGroupDO.Id, taskGroupDO.JobName, taskGroupDO.Caption, eumLogLevel.Warning, exp.Message)
		taskGroupDO.Cancel()
		taskGroupRepository.Save(taskGroupDO)
	})
}
