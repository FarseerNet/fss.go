package request

import (
	"github.com/farseer-go/fs/core/eumLogLevel"
	"time"
)

type LogRequest struct {
	LogLevel eumLogLevel.Enum // 日志等级
	Log      string           // 日志内容
	CreateAt time.Time        // 记录时间
}
