package taskGroup

import "time"

type DTO struct {
	// 主键
	Id int
	// 任务组标题
	Caption string
	// 实现Job的特性名称（客户端识别哪个实现类）
	JobName string
	// 传给客户端的参数，按逗号分隔
	Data map[string]string
	// 开始时间
	StartAt time.Time
	// 下次执行时间
	NextAt time.Time
	// 时间间隔
	IntervalMs int64
	// 时间定时器表达式
	Cron string
	// 活动时间
	ActivateAt time.Time
	// 最后一次完成时间
	LastRunAt time.Time
	// 是否开启
	IsEnable bool
	// 运行平均耗时
	RunSpeedAvg int64
	// 运行次数
	RunCount int
}
