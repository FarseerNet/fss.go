package job

import (
	"fss/application/tasks/taskGroupApp"
	"fss/domain/tasks/taskGroup"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"time"
)

func InitSysTaskJob() {
	flog.Info("正在读取所有任务组信息")
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	lstTaskGroup := taskGroupApp.ToList()
	flog.Infof("共获取到%d 条任务组信息", lstTaskGroup.Count())

	// 强制从缓存中再读一次，可以实现当缓存丢失时，可以重新写入该条任务组到缓存
	for _, taskGroupDO := range lstTaskGroup.ToArray() {
		taskGroupRepository.ToEntity(taskGroupDO.Id)
	}
}

func createSysJob(lstTaskGroup collections.List[taskGroup.DomainObject], jobName string, caption string, intervalMs time.Duration, data map[string]string) {
	do := lstTaskGroup.Where(func(item taskGroup.DomainObject) bool {
		return item.JobName == jobName
	}).First()

	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	if do.Id < 1 {
		if data == nil {
			data = make(map[string]string)
		}
		taskGroupDTO := taskGroupApp.DTO{
			Caption:    caption,
			JobName:    jobName,
			Data:       data,
			Cron:       parse.Convert(intervalMs.Milliseconds(), ""),
			StartAt:    time.Now(),
			NextAt:     time.Now(),
			ActivateAt: time.Now(),
			LastRunAt:  time.Now(),
			IsEnable:   true,
		}
		taskGroupDTO.Id = taskGroupApp.Add(taskGroupDTO)
		do = taskGroupRepository.ToEntity(taskGroupDTO.Id)
		lstTaskGroup.Add(do)
	} else if !do.IsEnable {
		do = taskGroupRepository.ToEntity(do.Id)
		do.SetEnable(true)
		taskGroupRepository.Save(do)
	}
}
