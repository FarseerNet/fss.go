package event

import (
	"fss/domain/clients/client"
	"github.com/farseer-go/eventBus"
)

// ClientOfflineEventName 事件名称
const ClientOfflineEventName = "ClientOffline"

// ClientOfflineEvent 客户端离线通知
type ClientOfflineEvent struct {
	Client client.DomainObject
}

// PublishEvent 发布事件
func (receiver ClientOfflineEvent) PublishEvent() {
	eventBus.PublishEvent(ClientOfflineEventName, receiver)
}
