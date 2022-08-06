package repository

import (
	"fss/domain/_/eumTaskType"
	"fss/domain/tasks/taskGroup"
	"fss/domain/tasks/taskGroup/vo"
	"fss/infrastructure/repository/model"
	"github.com/farseernet/farseer.go/cache/redis"
	"github.com/farseernet/farseer.go/core"
	"github.com/farseernet/farseer.go/core/container"
	"github.com/farseernet/farseer.go/data"
	"github.com/farseernet/farseer.go/exception"
	"github.com/farseernet/farseer.go/linq"
	"github.com/farseernet/farseer.go/mapper"
	"github.com/farseernet/farseer.go/utils/times"
	"time"
)

func RegisterTaskGroupRepository() {
	// 注册仓储
	_ = container.Register(func() taskGroup.Repository {
		repository := data.NewContext[taskGroupRepository]("default")
		repository.Client = redis.NewClient("default")
		return repository
	})
}

type taskGroupRepository struct {
	taskGroup data.TableSet[model.TaskGroupPO] `data:"name=task_group"`
	task      data.TableSet[model.TaskPO]      `data:"name=task"`
	Client    *redis.Client
}

func (repository taskGroupRepository) ToEntity(taskGroupId int) taskGroup.DomainObject {
	po := repository.taskGroup.Where("Id = ?", taskGroupId).ToEntity()
	do := mapper.Single[taskGroup.DomainObject](po)
	return do
}

func (repository taskGroupRepository) TodayFailCount() int64 {
	return repository.task.Where("Status = ? and CreateAt >= ?", eumTaskType.Fail, times.GetDate()).Count()
}

func (repository taskGroupRepository) ToTaskSpeedList(taskGroupId int) []int64 {
	lstPO := repository.task.Where("TaskGroupId = ? and Status = ?", taskGroupId, eumTaskType.Success).Desc("CreateAt").Select("RunSpeed").Limit(100).ToList()
	var lstIds []int64
	linq.From(lstPO).Select(&lstIds, func(item model.TaskPO) any {
		return item.RunSpeed
	})
	return lstIds
}

func (repository taskGroupRepository) ToList() []taskGroup.DomainObject {
	lstPO := repository.taskGroup.ToList()
	lstDO := mapper.Array[taskGroup.DomainObject](lstPO)
	return lstDO
}

func (repository taskGroupRepository) ToListByGroupId(groupId int, pageSize int, pageIndex int) core.PageList[vo.TaskDO] {
	page := repository.task.Where("TaskGroupId = ?", groupId).Select("Id", "Caption", "Progress", "Status", "StartAt", "CreateAt", "ClientIp", "RunSpeed", "RunAt").Desc("CreateAt").ToPageList(pageSize, pageIndex)
	return mapper.PageList[vo.TaskDO](page.List, page.RecordCount)
}

func (repository taskGroupRepository) ToListByClientId(clientId int64) []taskGroup.DomainObject {
	lst := repository.ToList()
	return linq.From(lst).Where(func(item taskGroup.DomainObject) bool {
		return item.Task.Client.Id == clientId && item.Task.StartAt.UnixMicro() < time.Now().UnixMicro()
	}).ToArray()
}

func (repository taskGroupRepository) GetTaskGroupCount() int64 {
	return repository.taskGroup.Count()
}

func (repository taskGroupRepository) ToFinishList(taskGroupId int, top int) []vo.TaskDO {
	lstPO := repository.task.Where("TaskGroupId = ? and (Status = ? or Status = ?)", taskGroupId, eumTaskType.Success, eumTaskType.Fail).Desc("CreateAt").Limit(top).ToList()
	lstDO := mapper.Array[vo.TaskDO](lstPO)
	return lstDO
}

func (repository taskGroupRepository) AddTask(taskDO vo.TaskDO) {
	po := mapper.Single[model.TaskPO](taskDO)
	repository.task.Insert(&po)
}

func (repository taskGroupRepository) Add(do taskGroup.DomainObject) taskGroup.DomainObject {
	po := mapper.Single[model.TaskGroupPO](do)
	repository.taskGroup.Insert(&po)
	do.Id = po.Id
	return do
}

func (repository taskGroupRepository) Save(do taskGroup.DomainObject) taskGroup.DomainObject {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) Delete(taskGroupId int) {
	repository.taskGroup.Where("Id = ?", taskGroupId).Delete()
	repository.task.Where("TaskGroupId = ?", taskGroupId).Delete()
}

func (repository taskGroupRepository) SyncToData() {
	//todo 这里需要读缓存
	lst := repository.ToList()
	for _, do := range lst {
		po := mapper.Single[model.TaskGroupPO](do)
		repository.taskGroup.Where("Id = ?", do.Id).Update(po)
	}
}

func (repository taskGroupRepository) GetCanSchedulerTaskGroup(jobsName []string, ts time.Duration, count int, client vo.ClientVO) []vo.TaskDO {
	getLocker := repository.Client.Lock.GetLocker("FSS_Scheduler", 5*time.Second)
	if !getLocker.TryLock() {
		exception.ThrowRefuseException("加锁失败")
	}
	defer getLocker.ReleaseLock()
	lstTaskGroup := repository.ToList()
	lstSchedulerTaskGroup := linq.From(lstTaskGroup).Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable &&
			linq.From(jobsName).ContainsItem(item.JobName) &&
			item.Task.Status == eumTaskType.None &&
			item.Task.StartAt.UnixMicro() < time.Now().Add(ts).UnixMicro() &&
			item.Task.Client.Id == 0
	}).OrderBy(func(item taskGroup.DomainObject) any {
		return item.StartAt.UnixMicro()
	}).Take(count)

	var lst []vo.TaskDO
	for _, taskGroupDO := range lstSchedulerTaskGroup {
		// 设为调度状态
		taskGroupDO.Scheduler(client)
		repository.Save(taskGroupDO)
		// 如果不相等，说明被其它客户端拿了
		if taskGroupDO.Task.Client.Id == client.Id {
			lst = append(lst, taskGroupDO.Task)
		}
	}
	return lst
}

func (repository taskGroupRepository) ToUnRunCount() int {
	lstTaskGroup := repository.ToList()
	return linq.From(lstTaskGroup).Where(func(item taskGroup.DomainObject) bool {
		return item.Task.Status == eumTaskType.None || item.Task.Status == eumTaskType.Scheduler || item.Task.CreateAt.UnixMicro() < time.Now().UnixMicro()
	}).Count()
}

func (repository taskGroupRepository) ToSchedulerWorkingList() []taskGroup.DomainObject {
	lstTaskGroup := repository.ToList()
	return linq.From(lstTaskGroup).Where(func(item taskGroup.DomainObject) bool {
		return item.Task.Status == eumTaskType.Scheduler || item.Task.Status == eumTaskType.Working
	}).ToArray()
}

func (repository taskGroupRepository) ToFinishPageList(pageSize int, pageIndex int) core.PageList[vo.TaskDO] {
	pageList := repository.task.Where("(Status = ? or Status = ?) and (CreateAt >= ?)", eumTaskType.Fail, eumTaskType.Success, time.Now().Add(-24*time.Hour)).
		Select("Id", "Caption", "Progress", "Status", "StartAt", "CreateAt", "ClientIp", "RunSpeed", "RunAt", "JobName").
		Desc("RunAt").ToPageList(pageSize, pageIndex)
	return mapper.PageList[vo.TaskDO](pageList.List, pageList.RecordCount)
}

func (repository taskGroupRepository) GetTaskUnFinishList(jobsName []string, top int) []taskGroup.DomainObject {
	lstTaskGroup := repository.ToList()
	return linq.From(lstTaskGroup).Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable && linq.From(jobsName).ContainsItem(item.JobName) && item.Task.Status != eumTaskType.Success && item.Task.Status != eumTaskType.Fail
	}).OrderBy(func(item taskGroup.DomainObject) any {
		return item.StartAt.UnixMicro()
	}).Take(top)
}

func (repository taskGroupRepository) GetEnableTaskList(status eumTaskType.Enum, pageSize int, pageIndex int) core.PageList[vo.TaskDO] {
	lstTaskGroup := repository.ToList()
	lstTaskGroup = linq.From(lstTaskGroup).Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable
	}).ToArray()

	if status != eumTaskType.None {
		lstTaskGroup = linq.From(lstTaskGroup).Where(func(item taskGroup.DomainObject) bool {
			return item.Task.Status == status
		}).ToArray()
	}

	lstTaskGroup = linq.From(lstTaskGroup).OrderBy(func(item taskGroup.DomainObject) any {
		return item.JobName
	}).ToArray()

	var lstTaskDO []vo.TaskDO
	linq.From(lstTaskGroup).Select(&lstTaskDO, func(item taskGroup.DomainObject) any {
		return item.Task
	})

	return linq.From(lstTaskDO).ToPageList(pageSize, pageIndex)
}
