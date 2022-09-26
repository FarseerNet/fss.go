package taskLog

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
)

type Repository interface {
	// GetList 获取日志
	GetList(jobName string, logLevel eumLogLevel.Enum, pageSize int, pageIndex int) collections.PageList[DomainObject]
	// Add 添加日志
	Add(taskLogDO DomainObject)
}
