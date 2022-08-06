package taskGroup

import (
	"fss/application/tasks/taskGroup/request"
	"fss/domain/log"
	"fss/domain/tasks"
	"fss/domain/tasks/taskGroup"
	"github.com/farseernet/farseer.go/core/container"
	"github.com/farseernet/farseer.go/core/eumLogLevel"
	"github.com/farseernet/farseer.go/exception"
	"github.com/farseernet/farseer.go/mapper"
)

type taskGroupApp struct {
	repository taskGroup.Repository
}

func NewTaskGroupApp() *taskGroupApp {
	return &taskGroupApp{repository: container.Resolve[taskGroup.Repository]()}
}

// CopyTaskGroup 复制任务组
func (r *taskGroupApp) CopyTaskGroup(request request.OnlyIdRequest) int {
	taskGroupDO := r.repository.ToEntity(request.Id)
	if taskGroupDO.Id < 1 {
		exception.ThrowRefuseException("要复制的任务组不存在")
	}

	newTaskGroup := taskGroup.Copy(taskGroupDO)
	r.repository.Add(newTaskGroup)
	return newTaskGroup.Id
}

// Delete Task 删除任务组
func (r *taskGroupApp) Delete(request request.OnlyIdRequest) {
	tasks.TaskGroupDeleteService(request.Id)
}

// Add 添加任务组信息
func (r *taskGroupApp) Add(dto DTO) int {
	if dto.Caption == "" || dto.Cron == "" || dto.Data == nil || dto.JobName == "" {
		exception.ThrowRefuseException("标题、时间间隔、传输数据、Job名称 必须填写")
	}
	do := mapper.Single[taskGroup.DomainObject](dto)
	do.CheckInterval()
	r.repository.Add(do)

	do.CreateTask()
	r.repository.Save(do)
	return do.Id
}

// Save 保存任务组
func (r *taskGroupApp) Save(dto DTO) {
	do := r.repository.ToEntity(dto.Id)
	if do.Id < 1 {
		exception.ThrowRefuseException("任务组不存在")
	}

	do.Set(dto.JobName, dto.Caption, dto.Data, dto.StartAt)
	do.SetCron(dto.Cron, dto.IntervalMs)
	do.SetEnable(dto.IsEnable)
	r.repository.Save(do)
}

// CancelTask 取消任务
func (r *taskGroupApp) CancelTask(request request.OnlyIdRequest) {
	do := r.repository.ToEntity(request.Id)
	do.Cancel()
	r.repository.Save(do)

	log.TaskLogAddService(request.Id, do.JobName, do.Caption, eumLogLevel.Information, "手动取消任务")
}

// SyncTaskGroup Task 同步数据
func (r *taskGroupApp) SyncTaskGroup() {
	r.repository.SyncToData()
}
