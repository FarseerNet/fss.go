package vo

import (
	"fss/domain/_/eumTaskType"
	"time"
)

// TaskEO 任务记录
type TaskEO struct {
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
	// 客户端
	Client ClientVO
	// 进度0-100
	Progress int
	// 状态
	Status eumTaskType.Enum
	// 任务创建时间
	CreateAt time.Time
	// 调度时间
	SchedulerAt time.Time
	// 本次执行任务时的Data数据
	Data map[string]string
}

func NewTaskDO() *TaskEO {
	return &TaskEO{}
}

// SetClient 调度时设置客户端
func (do *TaskEO) SetClient(client ClientVO) {
	do.Status = eumTaskType.Scheduler
	do.SchedulerAt = time.Now()
	do.Client = client
}

// SetJobName 更新了JobName，则要立即更新Task的JobName
func (do *TaskEO) SetJobName(jobName string) {
	do.JobName = jobName
}

// SetFail 设备为失败
func (do *TaskEO) SetFail() {
	do.Status = eumTaskType.Fail
}
