package tasks

import (
	"fss/domain/tasks/taskGroup"
	"fss/domain/tasks/taskGroup/vo"
	"github.com/farseer-go/fs/container"
)

// UpdateAvgSpeed 计算平均耗时
func UpdateAvgSpeed(taskGroupId int) {
	repository := container.Resolve[taskGroup.Repository]()

	var speedList = repository.ToTaskSpeedList(taskGroupId)
	var runSpeedAvg = vo.NewTaskSpeed(speedList).GetAvgSpeed()

	if runSpeedAvg > 0 {
		var do = repository.ToEntity(taskGroupId)
		do.RunSpeedAvg = runSpeedAvg
		repository.Save(do)
	}
}
