package job

import (
	"fss/application/tasks/taskGroupApp"
	"fss/domain/tasks"
	"github.com/farseer-go/fss"
)

// RegisterSyncTaskGroupAvgSpeedJob 计算任务组的平均耗时
func RegisterSyncTaskGroupAvgSpeedJob() {
	fss.RegisterJob("FSS.SyncTaskGroupAvgSpeed", syncTaskGroupAvgSpeedJob)
}

func syncTaskGroupAvgSpeedJob(context fss.IFssContext) bool {
	for _, taskGroupDO := range taskGroupApp.ToList().ToArray() {
		tasks.UpdateAvgSpeed(taskGroupDO.Id)
	}
	return true
}
