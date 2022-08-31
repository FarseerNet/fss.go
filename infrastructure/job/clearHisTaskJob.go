package job

import (
	"fss/domain/tasks/taskGroup/vo"
	"fss/infrastructure/repository"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fss"
)

// RegisterClearHisTaskJob 自动清除历史任务记录
func RegisterClearHisTaskJob() {
	fss.RegisterJob("FSS.ClearHisTask", clearHisTaskJob)
}

func clearHisTaskJob(context fss.IFssContext) bool {
	_reservedTaskCount := configure.GetInt("FSS.ReservedTaskCount")
	taskGroupRepository := repository.NewTaskGroupRepository()

	curIndex := 0
	result := 0
	lst := taskGroupRepository.ToList()
	for _, taskGroupDO := range lst.ToArray() {
		curIndex++
		lstTask := taskGroupRepository.ToFinishList(taskGroupDO.Id, _reservedTaskCount)
		if lstTask.Count() == 0 {
			continue
		}

		result += lstTask.Count()
		var taskId = lstTask.Min(func(item vo.TaskEO) any {
			return item.Id
		}).(int)

		// 清除历史记录
		taskGroupRepository.ClearFinish(taskGroupDO.Id, taskId)
		context.SetProgress(curIndex / lst.Count() * 100)
	}
	return true
}
