package request

import (
	"fss/domain/_/eumTaskType"
	"github.com/farseer-go/collections"
)

type JobInvokeRequest struct {
	TaskGroupId  int                                    // 任务组ID
	NextTimespan int64                                  // 下次执行时间
	Progress     int                                    // 当前进度
	Status       eumTaskType.Enum                       // 执行状态
	RunSpeed     int64                                  // 执行速度
	Log          LogRequest                             // 日志
	Data         collections.Dictionary[string, string] // 数据
}
