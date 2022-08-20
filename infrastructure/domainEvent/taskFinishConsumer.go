package domainEvent

import (
	"fss/domain/tasks/taskGroup"
	"fss/domain/tasks/taskGroup/event"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs/container"
)

func SubscribeTaskFinishEvent() {
	eventBus.Subscribe(event.TaskFinishEventName, taskFinishConsumer)
}

// TaskFinishConsumer 构建完成事件
func taskFinishConsumer(message any, _ eventBus.EventArgs) {
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	taskFinishEvent := message.(event.TaskFinishEvent)
	taskGroupRepository.AddTask(taskFinishEvent.Task)
}
