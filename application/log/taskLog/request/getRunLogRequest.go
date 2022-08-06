package request

import "github.com/farseernet/farseer.go/core/eumLogLevel"

type GetRunLogRequest struct {
	JobName   string
	LogLevel  eumLogLevel.Enum
	PageSize  int
	PageIndex int
}
