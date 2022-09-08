package taskGroup

import (
	"fss/domain/_/eumTaskType"
	"fss/domain/tasks/taskGroup/vo"
	"github.com/farseer-go/collections"
	"time"
)

type Repository interface {
	// ToEntity 获取任务组信息
	ToEntity(taskGroupId int) DomainObject
	// TodayFailCount 今日执行失败数量
	TodayFailCount() int64
	// ToTaskSpeedList 当前任务组下所有任务的执行速度
	ToTaskSpeedList(taskGroupId int) []int64
	// ToList 获取所有任务组中的任务
	ToList() collections.List[DomainObject]
	// ToListByGroupId 获取指定任务组的任务列表（FOPS）
	ToListByGroupId(groupId int, pageSize int, pageIndex int) collections.PageList[vo.TaskEO]
	// ToListByClientId 获取指定客户端的任务组列表
	ToListByClientId(clientId int64) collections.List[DomainObject]
	// GetTaskGroupCount 获取任务组数量
	GetTaskGroupCount() int64
	// ToFinishList 获取指定任务组执行成功的任务列表
	ToFinishList(taskGroupId int, top int) collections.List[vo.TaskEO]
	// AddTask 创建任务
	AddTask(taskDO vo.TaskEO)
	// Add 添加任务组
	Add(do *DomainObject)
	// Save 保存任务组信息
	Save(do DomainObject)
	// Delete 删除任务组
	Delete(taskGroupId int)
	// GetCanSchedulerTaskGroup 获取所有任务组中的任务
	GetCanSchedulerTaskGroup(jobsName []string, ts time.Duration, count int, client vo.ClientVO) collections.List[vo.TaskEO]
	// ToUnRunCount 获取未执行的任务数量
	ToUnRunCount() int
	// ToSchedulerWorkingList 获取执行中的任务
	ToSchedulerWorkingList() collections.List[DomainObject]
	// ToFinishPageList 获取已完成的任务列表
	ToFinishPageList(pageSize int, pageIndex int) collections.PageList[vo.TaskEO]
	// GetTaskUnFinishList 获取进行中的任务
	GetTaskUnFinishList(jobsName []string, top int) collections.List[DomainObject]
	// GetEnableTaskList 获取在用的任务组
	GetEnableTaskList(status eumTaskType.Enum, pageSize int, pageIndex int) collections.PageList[vo.TaskEO]
	// ToIdList 从数据库中读取数据
	ToIdList() []int
	// SaveToDb 保存到数据库
	SaveToDb(do DomainObject)
	// ClearFinish 清除成功的任务记录（1天前）
	ClearFinish(groupId int, taskId int)
}
