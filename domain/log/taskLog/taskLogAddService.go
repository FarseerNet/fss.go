package taskLog

import (
	"github.com/farseernet/farseer.go/core/container"
	"github.com/farseernet/farseer.go/core/eumLogLevel"
	"time"
)

// TaskLogAddService 添加日志记录
func TaskLogAddService(taskGroupId int, jobName string, caption string, logLevel eumLogLevel.Enum, content string) {
	repository := container.Resolve[Repository]()
	repository.Add(DomainObject{
		TaskGroupId: taskGroupId,
		Caption:     caption,
		JobName:     jobName,
		LogLevel:    logLevel,
		Content:     content,
		CreateAt:    time.Now(),
	})
}
