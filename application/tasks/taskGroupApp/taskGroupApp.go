package taskGroupApp

import (
	"fss/application/tasks/taskGroupApp/request"
	"fss/domain/log"
	"fss/domain/tasks"
	"fss/domain/tasks/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
)

// CopyTaskGroup 复制任务组
func CopyTaskGroup(request request.OnlyIdRequest) int {
	repository := container.Resolve[taskGroup.Repository]()
	taskGroupDO := repository.ToEntity(request.Id)
	if taskGroupDO.IsNull() {
		exception.ThrowRefuseException("要复制的任务组不存在")
	}

	newTaskGroup := taskGroup.Copy(taskGroupDO)
	repository.Add(&newTaskGroup)
	return newTaskGroup.Id
}

// Delete Task 删除任务组
func Delete(request request.OnlyIdRequest) {
	tasks.DeleteTaskGroupService(request.Id)
}

// Add 添加任务组信息
func Add(dto DTO) int {
	repository := container.Resolve[taskGroup.Repository]()
	if dto.Caption == "" || dto.Cron == "" || dto.JobName == "" {
		exception.ThrowRefuseException("标题、时间间隔、传输数据、Job名称 必须填写")
	}
	do := mapper.Single[taskGroup.DomainObject](&dto)
	do.CheckInterval()
	repository.Add(&do)

	do.CreateTask()
	repository.Save(do)
	return do.Id
}

// Save 保存任务组
func Save(dto DTO) {
	repository := container.Resolve[taskGroup.Repository]()
	do := repository.ToEntity(dto.Id)
	if do.IsNull() {
		exception.ThrowRefuseException("任务组不存在")
	}

	do.Set(dto.JobName, dto.Caption, dto.Data, dto.StartAt)
	do.SetCron(dto.Cron, dto.IntervalMs)
	do.SetEnable(dto.IsEnable)
	repository.Save(do)
}

// CancelTask 取消任务
func CancelTask(request request.OnlyIdRequest) {
	repository := container.Resolve[taskGroup.Repository]()
	do := repository.ToEntity(request.Id)
	if !do.IsNull() {
		do.Cancel()
		repository.Save(do)
	}
	log.TaskLogAddService(request.Id, do.JobName, do.Caption, eumLogLevel.Information, "手动取消任务")
}

// SyncTaskGroup 同步数据
func SyncTaskGroup() {
	tasks.SyncTaskGroupService()
}

// SetEnable 设置任务组状态
func SetEnable(request request.SetEnableRequest) {
	repository := container.Resolve[taskGroup.Repository]()
	do := repository.ToEntity(request.Id)
	if do.IsNull() {
		exception.ThrowRefuseException("任务组不存在")
	}

	do.SetEnable(request.IsEnable)
	repository.Save(do)
}
