package request

import "github.com/farseer-go/fs/core/eumLogLevel"

type GetRunLogRequest struct {
	JobName   string
	LogLevel  eumLogLevel.Enum
	PageSize  int
	PageIndex int
}
