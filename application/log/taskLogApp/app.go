package taskLogApp

import (
	"fss/application/log/taskLogApp/request"
	"fss/domain/log"
	"fss/domain/log/taskLog"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
)

// Add 添加日志记录
func Add(taskGroupId int, jobName string, caption string, logLevel eumLogLevel.Enum, content string) {
	log.TaskLogAddService(taskGroupId, jobName, caption, logLevel, content)
}

// GetList 获取日志
func GetList(request request.GetRunLogRequest) collections.List[taskLog.DomainObject] {
	repository := container.Resolve[taskLog.Repository]()
	return repository.GetList(request.JobName, request.LogLevel, request.PageSize, request.PageIndex)
}
