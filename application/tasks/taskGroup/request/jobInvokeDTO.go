package request

import "fss/domain/_/eumTaskType"

type JobInvokeDTO struct {
	// 任务组ID
	TaskGroupId int
	// 下次执行时间
	NextTimespan int64
	// 当前进度
	Progress int
	// 执行状态
	Status eumTaskType.Enum
	// 执行速度
	RunSpeed int64
	// 日志
	Log LogRequest
	// 数据
	Data map[string]string
}
