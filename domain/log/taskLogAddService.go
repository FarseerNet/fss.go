package log

import (
	"fss/domain/log/taskLog"
	"github.com/farseer-go/fs/core/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"time"
)

// TaskLogAddService 添加日志记录
func TaskLogAddService(taskGroupId int, jobName string, caption string, logLevel eumLogLevel.Enum, content string) {
	repository := container.Resolve[taskLog.Repository]()
	repository.Add(taskLog.DomainObject{
		TaskGroupId: taskGroupId,
		Caption:     caption,
		JobName:     jobName,
		LogLevel:    logLevel,
		Content:     content,
		CreateAt:    time.Now(),
	})
}
