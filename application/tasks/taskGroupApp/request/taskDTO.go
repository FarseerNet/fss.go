package request

import (
	"fss/domain/_/eumTaskType"
	"time"
)

type TaskDTO struct {
	// 主键
	Id int
	// 任务组ID
	TaskGroupId int
	// 任务组标题
	Caption string
	// 实现Job的特性名称（客户端识别哪个实现类）
	JobName string
	// 开始时间
	StartAt time.Time
	// 实际执行时间
	RunAt time.Time
	// 运行耗时
	RunSpeed int64
	// 客户端Id
	ClientId int64
	// 客户端IP
	ClientIp string
	// 客户端名称
	ClientName string
	// 进度0-100
	Progress int
	// 状态
	Status eumTaskType.Enum
	// 任务创建时间
	CreateAt time.Time
	// 调度时间
	SchedulerAt time.Time
	// 数据
	Data map[string]string
}
