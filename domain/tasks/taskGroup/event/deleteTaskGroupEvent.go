package event

import "github.com/farseernet/farseer.go/eventBus"

// DeleteTaskGroupEventName 事件名称
const DeleteTaskGroupEventName = "DeleteTaskGroup"

// DeleteTaskGroupEvent 删除任务组事件
type DeleteTaskGroupEvent struct {
	TaskGroupId int
}

// PublishEvent 发布事件
func (receiver DeleteTaskGroupEvent) PublishEvent() {
	eventBus.PublishEvent(DeleteTaskGroupEventName, receiver)
}
