package domainEvent

import (
	"fss/domain/tasks/taskGroup"
	"fss/domain/tasks/taskGroup/event"
	"github.com/farseernet/farseer.go/core/container"
	"github.com/farseernet/farseer.go/eventBus"
)

func SubscribeTaskFinishEvent() {
	eventBus.Subscribe(event.TaskFinishEventName, taskFinishConsumer)
}

var taskGroupRepository = container.Resolve[taskGroup.Repository]()

// TaskFinishConsumer 构建完成事件
func taskFinishConsumer(message any, _ eventBus.EventArgs) {
	taskFinishEvent := message.(event.TaskFinishEvent)
	taskGroupRepository.AddTask(taskFinishEvent.Task)
}
