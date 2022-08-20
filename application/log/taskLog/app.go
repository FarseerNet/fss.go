package taskLog

import (
	"fss/application/log/taskLog/request"
	"fss/domain/log"
	"fss/domain/log/taskLog"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/mapper"
)

type app struct {
	repository taskLog.Repository
}

func NewApp() *app {
	return &app{repository: container.Resolve[taskLog.Repository]()}
}

// Add 添加日志记录
func (r *app) Add(taskGroupId int, jobName string, caption string, logLevel eumLogLevel.Enum, content string) {
	log.TaskLogAddService(taskGroupId, jobName, caption, logLevel, content)
}

// GetList 获取日志
func (r *app) GetList(request request.GetRunLogRequest) []Dto {
	lstDO := r.repository.GetList(request.JobName, request.LogLevel, request.PageSize, request.PageIndex)
	return mapper.Array[Dto](lstDO)
}
