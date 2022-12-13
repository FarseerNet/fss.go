package job

import (
	"fss/application/tasks/taskGroupApp"
	"fss/domain/tasks"
	"github.com/farseer-go/fss"
)

// SyncTaskGroupAvgSpeedJob 计算任务组的平均耗时
func SyncTaskGroupAvgSpeedJob(context fss.IFssContext) bool {
	for _, taskGroupDO := range taskGroupApp.ToList().ToArray() {
		tasks.UpdateAvgSpeed(taskGroupDO.Id)
	}
	return true
}
