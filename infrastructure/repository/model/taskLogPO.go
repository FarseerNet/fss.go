package model

import (
	"github.com/farseer-go/fs/core/eumLogLevel"
	"time"
)

type TaskLogPO struct {
	// 主键
	Id int64 `gorm:"primaryKey" es_type:"long"`
	// 任务组记录ID
	TaskGroupId int `es_type:"integer"`
	// 任务组标题
	Caption string `es_type:"keyword"`
	// 实现Job的特性名称（客户端识别哪个实现类）
	JobName string `es_type:"keyword"`
	// 日志级别
	LogLevel eumLogLevel.Enum `es_type:"byte"`
	// 日志内容
	Content string `es_type:"test"`
	// 日志时间
	CreateAt time.Time `es_type:"date"`
}
