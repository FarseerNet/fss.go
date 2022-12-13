package job

import (
	"fss/domain/tasks"
	"github.com/farseer-go/fss"
)

// SyncTaskGroupJob 同步任务组信息数据库与缓存
func SyncTaskGroupJob(context fss.IFssContext) bool {
	tasks.SyncTaskGroupService()
	return true
}
