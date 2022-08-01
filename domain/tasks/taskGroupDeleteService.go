package tasks

import (
	"fss/domain/tasks/taskGroup"
	"fss/domain/tasks/taskGroup/event"
)

type TaskGroupDeleteService struct {
	repository taskGroup.Repository
}

// Delete 删除任务组
func (service TaskGroupDeleteService) Delete(taskGroupId int) {
	var do = service.repository.ToEntity(taskGroupId)
	if do.Id < 1 {
		panic("要删除的任务组不存在")
	}

	// 如果任务组是开启状态，则需要先暂定任务组
	if do.IsEnable {
		do.Disable()
		service.repository.Save(do)
	}

	service.repository.Delete(do.Id)

	// 发布删除任务组事件
	event.DeleteTaskGroupEvent{TaskGroupId: do.Id}.PublishEvent()
}
