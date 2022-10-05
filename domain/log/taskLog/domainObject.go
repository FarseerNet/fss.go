package taskLog

import (
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"log"
	"time"
)

type DomainObject struct {
	TaskGroupId int              // 任务组记录ID
	Caption     string           // 任务组标题
	JobName     string           // 实现Job的特性名称（客户端识别哪个实现类）
	LogLevel    eumLogLevel.Enum // 日志级别
	Content     string           // 日志内容
	CreateAt    time.Time        // 日志时间
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
