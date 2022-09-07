package taskGroup

import (
	"fss/domain/_/eumTaskType"
	"fss/domain/tasks/taskGroup/event"
	"fss/domain/tasks/taskGroup/vo"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/utils/times"
	"github.com/robfig/cron/v3"
	"math"
	"strconv"
	"strings"
	"time"
)

type DomainObject struct {
	// 主键
	Id int
	// 任务
	Task vo.TaskEO
	// 任务组标题
	Caption string
	// 实现Job的特性名称（客户端识别哪个实现类）
	JobName string
	// 本次执行任务时的Data数据
	Data collections.Dictionary[string, string]
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

// Copy 复制新的任务组
func Copy(copySource DomainObject) DomainObject {
	return DomainObject{
		Caption:    copySource.Caption + "复制",
		JobName:    copySource.JobName,
		Data:       copySource.Data,
		IntervalMs: copySource.IntervalMs,
		Cron:       copySource.Cron,
		IsEnable:   copySource.IsEnable,
	}
}

// CheckInterval 保存任务组信息前，检查状态
func (do *DomainObject) CheckInterval() {
	// 是否为数字
	if parse.IsInt(do.Cron) {
		do.IntervalMs = parse.Convert(do.Cron, int64(0))
		do.Cron = ""
		do.NextAt = time.Now().Add(time.Duration(do.IntervalMs) * time.Millisecond)
	} else {
		cornSchedule, err := cron.ParseStandard(do.Cron)
		if err != nil {
			exception.ThrowRefuseException("Cron格式错误")
		}
		do.IntervalMs = 0
		do.NextAt = cornSchedule.Next(time.Now())
	}

	do.Caption = strings.Trim(do.Caption, " ")
	do.JobName = strings.Trim(do.JobName, " ")
	do.Cron = strings.Trim(do.Cron, " ")
}

// Disable 修改任务状态为不可用
func (do *DomainObject) Disable() {
	do.IsEnable = false
}

// Set 修改了任务组信息
func (do *DomainObject) Set(jobName string, caption string, data collections.Dictionary[string, string], startAt time.Time) {
	// 更新了JobName，则要立即更新Task的JobName
	if do.JobName != jobName && do.Task.Status == eumTaskType.None {
		do.Task.SetJobName(jobName)
	}
	do.JobName = jobName
	do.Caption = caption
	do.Data = data
	do.StartAt = startAt
}

// SetCron 修改了任务的时间间隔
func (do *DomainObject) SetCron(strCron string, intervalMs int64) {
	// 是否为数字
	if do.Cron != strCron || do.IntervalMs != intervalMs {
		// 是否为数字
		if parse.IsInt(do.Cron) {
			do.IntervalMs = parse.Convert(do.Cron, int64(0))
			do.Cron = ""
			do.NextAt = time.Now().Add(time.Duration(do.IntervalMs) * time.Millisecond)
		} else {
			cornSchedule, err := cron.ParseStandard(do.Cron)
			if err != nil {
				exception.ThrowRefuseException("Cron格式错误")
			}
			do.IntervalMs = 0
			do.NextAt = cornSchedule.Next(time.Now())
		}
	}
}

// SetEnable 更改启用状态
func (do *DomainObject) SetEnable(enable bool) {
	// 停止了任务，需要把任务取消掉
	if do.IsEnable && !enable {
		do.IsEnable = enable
		do.Cancel()
	} else if !do.IsEnable && enable { // 重新开启了任务
		do.IsEnable = enable
		switch do.Task.Status {
		// 进行中的任务，要先取消
		case eumTaskType.Scheduler:
		case eumTaskType.Working:
			do.Cancel()
			break
			// 未开始的任务，直接保存
		case eumTaskType.None:
		case eumTaskType.Fail:
		case eumTaskType.Success:
			do.Finish()
			break
		}
	}
}

// Cancel 取消任务
func (do *DomainObject) Cancel() {
	do.Task.SetFail()
	// 这里不设置的话，下次执行时间，有可能还是将来的，导致如果设置错了时间的话，会一直等待原来设置错的时间
	do.NextAt = do.LastRunAt
	// 设置下一次的执行时间
	do.CalculateNextTime()
	// 创建新的任务
	do.CreateTask()
}

// Finish 保存Task（taskGroup必须是最新的）
func (do *DomainObject) Finish() {
	do.CalculateNextTime()
	// 如果是停止状态，创建任务不会执行。则需要在这里进行保存
	do.CreateTask()
}

// CalculateNextTime 任务完成后，计算下一次的时间
func (do *DomainObject) CalculateNextTime() {
	// 本次的时间策略晚，则通过时间策略计算出来
	if time.Now().UnixMicro() > do.NextAt.UnixMicro() {

		// 时间间隔器
		if do.IntervalMs > 0 {
			do.NextAt = time.Now().Add(time.Duration(do.IntervalMs) * time.Millisecond)
		} else {
			cornSchedule, err := cron.ParseStandard(do.Cron)
			if do.Cron != "" && err != nil {
				do.NextAt = cornSchedule.Next(time.Now())
			} else { // 没有找到设置下一次时间的设置，则默认30S执行一次
				do.NextAt = time.Now().Add(time.Duration(30) * time.Second)
			}
		}
	}
}

// CreateTask 创建新的Task
func (do *DomainObject) CreateTask() {
	if !do.IsEnable {
		do.Task = vo.TaskEO{}
		return
	}

	if do.Task.IsFinish() {
		// 任务完成，发布完成事件
		event.TaskFinishEvent{Task: do.Task}.PublishEvent()
	}

	if do.Task.IsNull() || do.Task.IsFinish() {
		// 没查到时，自动创建一条对应的Task
		do.Task = vo.TaskEO{
			TaskGroupId: do.Id,
			StartAt:     do.NextAt,
			Caption:     do.Caption,
			JobName:     do.JobName,
			RunSpeed:    0,
			Client:      vo.ClientVO{},
			Progress:    0,
			Status:      eumTaskType.None,
			CreateAt:    time.Now(),
			RunAt:       time.Now(),
			SchedulerAt: time.Now(),
			Data:        do.Data,
		}
	}
}

// Scheduler 调度时设置客户端
func (do *DomainObject) Scheduler(client vo.ClientVO) {
	if do.Task.Status == eumTaskType.None {
		do.Task.SetClient(client)
		do.Task.Data = do.Data
		do.ActivateAt = time.Now()
	}
}

// CheckClientOffline 检测进行中状态的任务
func (do *DomainObject) CheckClientOffline() {
	if do.Task.Status != eumTaskType.Scheduler && do.Task.Status != eumTaskType.Working {
		return
	}

	// 任务组活动时间大于1分钟，判定为客户端下线
	if time.Now().Unix()-do.ActivateAt.Unix() >= 60 { // 大于1分钟，才检查
		exception.ThrowRefuseException("【任务假死】任务：【" + do.JobName + "】 " + strconv.Itoa(do.Id) + " " + do.Caption + " " + do.Task.Status.String() + " 在" + times.GetSubDesc(time.Now(), do.ActivateAt) + "没有反应，强制设为失败状态")
	}

	// 如果时间小于5分钟的，则按5分钟来判定
	var timeout = math.Max(float64(do.RunSpeedAvg)*2.5, float64((5 * time.Minute).Milliseconds()))

	if float64(time.Now().Sub(do.Task.RunAt).Milliseconds()) > timeout {
		exception.ThrowRefuseException("【任务超时】任务：【" + do.JobName + "】 " + strconv.Itoa(do.Id) + " " + do.Caption + " 超过平均运行时间：" + parse.Convert(timeout, "0") + " ms，强制设为失败状态")
	}
}

// Working 执行中
func (do *DomainObject) Working(data collections.Dictionary[string, string], nextTimespan int64, progress int, status eumTaskType.Enum, runSpeed int64) {
	// 数据库的状态处于调度状态，说明客户端第一次请求进来
	if do.Task.Status == eumTaskType.Scheduler {
		do.Task.RunAt = time.Now() // 首次执行，记录时间
		// 更新group元信息
		do.RunCount++
		do.LastRunAt = time.Now()
	}

	do.Data = data
	do.ActivateAt = time.Now()
	do.Task.Progress = progress
	do.Task.Status = status
	do.Task.RunSpeed = runSpeed

	// 客户端设置了动态时间
	if nextTimespan > 0 {
		do.NextAt = time.UnixMilli(nextTimespan)
	}

	// 如果=成功|错误状态，则要立即更新数据库
	switch do.Task.Status {
	case eumTaskType.Fail, eumTaskType.Success:
		do.Finish()
		break
	}
}
