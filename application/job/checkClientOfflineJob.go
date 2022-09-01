package job

import (
	"fss/domain/clients/client"
	"fss/domain/clients/client/event"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fss"
)

// RegisterCheckClientOfflineJob 检查超时离线的客户端
func RegisterCheckClientOfflineJob() {
	fss.RegisterJob("FSS.CheckClientOffline", checkClientOfflineJob)
}

func checkClientOfflineJob(context fss.IFssContext) bool {
	clientRepository := container.Resolve[client.Repository]()
	// 拿到超时的客户端列表
	lstTimeout := clientRepository.ToList().Where(func(item client.DomainObject) bool {
		return item.IsTimeout()
	}).ToArray()

	for _, clientDO := range lstTimeout {
		clientRepository.RemoveClient(clientDO.Id)
		event.ClientOfflineEvent{Client: clientDO}.PublishEvent()
	}
	return true
}
