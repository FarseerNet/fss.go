package job

import (
	"fss/application/tasks/taskGroupApp"
	"fss/domain/tasks/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
	"time"
)

// CheckFinishStatusJob 检测完成状态的任务
func CheckFinishStatusJob(context *tasks.TaskContext) {
	dicTaskGroup := taskGroupApp.ToList()
	var ids []int
	dicTaskGroup.Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable
	}).Select(&ids, func(item taskGroup.DomainObject) any {
		return item.Id
	})

	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	for _, taskGroupId := range ids {
		taskGroupDO := taskGroupRepository.ToEntity(taskGroupId)
		if taskGroupDO.IsNull() {
			continue
		}
		// 加个时间，来限制并发
		if time.Now().Sub(taskGroupDO.Task.RunAt).Seconds() < 30 {
			continue
		}

		if taskGroupDO.Task.IsFinish() || taskGroupDO.Task.IsNull() {
			taskGroupDO.CreateTask()
			taskGroupRepository.Save(taskGroupDO)
		}
	}
}
