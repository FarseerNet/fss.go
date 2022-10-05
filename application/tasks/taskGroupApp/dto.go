package taskGroupApp

import (
	"github.com/farseer-go/collections"
	"time"
)

type DTO struct {
	Id          int                                    // 主键
	Caption     string                                 // 任务组标题
	JobName     string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Data        collections.Dictionary[string, string] // 传给客户端的参数，按逗号分隔
	StartAt     time.Time                              // 开始时间
	NextAt      time.Time                              // 下次执行时间
	IntervalMs  int64                                  // 时间间隔
	Cron        string                                 // 时间定时器表达式
	ActivateAt  time.Time                              // 活动时间
	LastRunAt   time.Time                              // 最后一次完成时间
	IsEnable    bool                                   // 是否开启
	RunSpeedAvg int64                                  // 运行平均耗时
	RunCount    int                                    // 运行次数
}
