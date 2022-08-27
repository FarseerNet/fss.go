package request

import (
	"github.com/farseer-go/fs/core/eumLogLevel"
	"time"
)

type LogRequest struct {
	// 日志等级
	LogLevel eumLogLevel.Enum
	// 日志内容
	Log string
	// 记录时间
	CreateAt time.Time
}
