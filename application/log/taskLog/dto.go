package taskLog

import (
	"github.com/farseer-go/fs/core/eumLogLevel"
	"time"
)

type Dto struct {
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
