package taskLog

import (
	"github.com/farseernet/farseer.go/core/container"
	"github.com/farseernet/farseer.go/core/eumLogLevel"
	"log"
	"time"
)

type DomainObject struct {
	// 主键
	Id int64
	// 任务组记录ID
	TaskGroupId int
	// 任务组标题
	Caption string
	// 实现Job的特性名称（客户端识别哪个实现类）
	JobName string
	// 日志级别
	LogLevel eumLogLevel.Enum
	// 日志内容
	Content string
	// 日志时间
	CreateAt time.Time
}

func NewDO(taskGroupId int, jobName string, caption string, logLevel eumLogLevel.Enum, content string) DomainObject {
	do := DomainObject{
		TaskGroupId: taskGroupId,
		Caption:     caption,
		JobName:     jobName,
		LogLevel:    logLevel,
		Content:     content,
		CreateAt:    time.Now(),
	}
	if logLevel == eumLogLevel.Error || logLevel == eumLogLevel.Warning {
		log.Println(eumLogLevel.GetName(logLevel) + "：" + content)
	}
	return do
}

// Add 添加日志到队列
func (do DomainObject) Add() {
	repository := container.Resolve[Repository]()
	repository.Add(do)
}
