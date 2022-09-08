package tasks

import (
	"fss/domain/tasks/taskGroup"
	"fss/domain/tasks/taskGroup/event"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
)

// DeleteTaskGroupService 删除任务组
func DeleteTaskGroupService(taskGroupId int) {
	repository := container.Resolve[taskGroup.Repository]()

	var do = repository.ToEntity(taskGroupId)
	if do.Id < 1 {
		exception.ThrowRefuseException("要删除的任务组不存在")
	}

	// 如果任务组是开启状态，则需要先暂定任务组
	if do.IsEnable {
		do.Disable()
		repository.Save(do)
	}

	repository.Delete(do.Id)

	// 发布删除任务组事件
	event.DeleteTaskGroupEvent{TaskGroupId: do.Id}.PublishEvent()
}
