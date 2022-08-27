package taskLogApp

import (
	"fss/application/log/taskLogApp/request"
	"fss/domain/log"
	"fss/domain/log/taskLog"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/mapper"
)

// Add 添加日志记录
func Add(taskGroupId int, jobName string, caption string, logLevel eumLogLevel.Enum, content string) {
	log.TaskLogAddService(taskGroupId, jobName, caption, logLevel, content)
}

// GetList 获取日志
func GetList(request request.GetRunLogRequest) []Dto {
	repository := container.Resolve[taskLog.Repository]()
	lstDO := repository.GetList(request.JobName, request.LogLevel, request.PageSize, request.PageIndex)
	return mapper.Array[Dto](lstDO)
}
