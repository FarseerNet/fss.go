package job

import (
	"fss/infrastructure/repository"
	"github.com/farseer-go/fss"
)

// RegisterSyncTaskGroupJob 同步任务组信息数据库与缓存
func RegisterSyncTaskGroupJob() {
	fss.RegisterJob("FSS.SyncTaskGroup", syncTaskGroupJob)
}

func syncTaskGroupJob(context fss.IFssContext) bool {
	taskGroupRepository := repository.NewTaskGroupRepository()
	lstGroupByDb := taskGroupRepository.ToDbList()
	curIndex := 0

	for _, taskGroupVo := range lstGroupByDb.ToArray() {
		curIndex++
		// 强制从缓存中再读一次，可以实现当缓存丢失时，可以重新写入该条任务组到缓存
		po := taskGroupRepository.ToEntity(taskGroupVo.Id)
		taskGroupRepository.Save(po)
		context.SetProgress(curIndex / lstGroupByDb.Count() * 100)
	}
	return true
}
