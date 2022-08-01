package tasks

import (
	"fss/domain/tasks/taskGroup"
	"fss/domain/tasks/taskGroup/vo"
)

type TaskGroupService struct {
	repository taskGroup.Repository
}

// UpdateAvgSpeed 计算平均耗时
func (service TaskGroupService) UpdateAvgSpeed(taskGroupId int) {
	var speedList = service.repository.ToTaskSpeedList(taskGroupId)
	var runSpeedAvg = vo.NewTaskSpeed(speedList).GetAvgSpeed()

	if runSpeedAvg > 0 {
		var do = service.repository.ToEntity(taskGroupId)
		do.RunSpeedAvg = runSpeedAvg
		service.repository.Save(do)
	}
}
