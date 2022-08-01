package taskLog

import "github.com/farseernet/farseer.go/core/eumLogLevel"

type Repository interface {
	// GetList 获取日志
	GetList(jobName string, logLevel eumLogLevel.Enum, pageSize int, pageIndex int) []DomainObject
	// Add 添加日志
	Add(taskLogDO DomainObject)
}
