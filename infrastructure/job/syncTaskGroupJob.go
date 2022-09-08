package job

import (
	"fss/domain/tasks"
	"github.com/farseer-go/fss"
)

// RegisterSyncTaskGroupJob 同步任务组信息数据库与缓存
func RegisterSyncTaskGroupJob() {
	fss.RegisterJob("FSS.SyncTaskGroup", syncTaskGroupJob)
}

func syncTaskGroupJob(context fss.IFssContext) bool {
	tasks.SyncTaskGroupService()
	return true
}
