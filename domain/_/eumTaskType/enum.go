package eumTaskType

type Enum int

const (
	// None 未开始
	None Enum = iota
	// Scheduler 已调度
	Scheduler
	// Working 执行中
	Working
	// Fail 失败
	Fail
	// Success 完成
	Success
)

func (e Enum) String() string {
	switch e {
	case None:
		return "None"
	case Scheduler:
		return "Scheduler"
	case Working:
		return "Working"
	case Fail:
		return "Fail"
	case Success:
		return "Success"
	}
	return "None"
}
