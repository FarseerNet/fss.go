package tasks

import (
	"fss/domain/tasks/taskGroup"
	"github.com/farseer-go/fs/container"
)

// SyncTaskGroupService 同步任务组信息数据库与缓存
func SyncTaskGroupService() {
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	lstIds := taskGroupRepository.ToIdList()

	for _, id := range lstIds {
		// 从缓存中读取，然后写入数据库
		po := taskGroupRepository.ToEntity(id)
		if !po.IsNull() {
			taskGroupRepository.SaveToDb(po)
		}
	}
}
