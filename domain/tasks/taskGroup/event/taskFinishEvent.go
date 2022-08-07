package event

import (
	"fss/domain/tasks/taskGroup/vo"
	"github.com/farseer-go/eventBus"
)

// TaskFinishEventName 事件名称
const TaskFinishEventName = "TaskFinish"

// TaskFinishEvent 任务完成事件
type TaskFinishEvent struct {
	Task vo.TaskDO
}

// PublishEvent 发布事件
func (receiver TaskFinishEvent) PublishEvent() {
	eventBus.PublishEvent(TaskFinishEventName, receiver)
}
