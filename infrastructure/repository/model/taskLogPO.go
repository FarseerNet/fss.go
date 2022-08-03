package model

import (
	"github.com/farseernet/farseer.go/core/eumLogLevel"
	"time"
)

type TaskLogPO struct {
	// 主键
	//[Number(type: NumberType.Long)]
	Id int64 `gorm:"primaryKey"`
	// 任务组记录ID
	//[Number(type: NumberType.Integer)]
	TaskGroupId int
	// 任务组标题
	//[Keyword]
	Caption string
	// 实现Job的特性名称（客户端识别哪个实现类）
	// [Keyword]
	JobName string
	// 日志级别
	//[Number(type: NumberType.Byte)]
	LogLevel eumLogLevel.Enum
	// 日志内容
	//[Text]
	Content string
	// 日志时间
	// [Date]
	CreateAt time.Time
}
