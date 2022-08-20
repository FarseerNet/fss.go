package repository

import (
	"fss/domain/_/eumTaskType"
	"fss/domain/tasks/taskGroup"
	"fss/domain/tasks/taskGroup/vo"
	"fss/infrastructure/repository/model"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/redis"
	"github.com/farseer-go/utils/times"
	"strconv"
	"time"
)

const taskGroupCacheKey = "FSS_TaskGroup"

func RegisterTaskGroupRepository() {
	// 注册仓储
	container.Register(func() taskGroup.Repository {
		repository := data.NewContext[taskGroupRepository]("default")
		repository.redis = redis.NewClient("default")
		return repository
	})
}

type taskGroupRepository struct {
	taskGroup data.TableSet[model.TaskGroupPO] `data:"name=task_group"`
	task      data.TableSet[model.TaskPO]      `data:"name=task"`
	redis     *redis.Client
}

func (repository taskGroupRepository) ToEntity(taskGroupId int) taskGroup.DomainObject {
	var do taskGroup.DomainObject
	if repository.redis.Hash.ToEntity(taskGroupCacheKey, strconv.Itoa(taskGroupId), &do) == nil {
		po := repository.taskGroup.Where("Id = ?", taskGroupId).ToEntity()
		do = mapper.Single[taskGroup.DomainObject](po)
	}
	return do
}

func (repository taskGroupRepository) TodayFailCount() int64 {
	return repository.task.Where("Status = ? and CreateAt >= ?", eumTaskType.Fail, times.GetDate()).Count()
}

func (repository taskGroupRepository) ToTaskSpeedList(taskGroupId int) []int64 {
	lstPO := repository.task.Where("TaskGroupId = ? and Status = ?", taskGroupId, eumTaskType.Success).Desc("CreateAt").Select("RunSpeed").Limit(100).ToList()
	var lstIds []int64
	lstPO.Select(&lstIds, func(item model.TaskPO) any {
		return item.RunSpeed
	})
	return lstIds
}

func (repository taskGroupRepository) ToList() collections.List[taskGroup.DomainObject] {
	var lstDO collections.List[taskGroup.DomainObject]
	repository.taskGroup.ToList().MapToList(&lstDO)
	return lstDO
}

func (repository taskGroupRepository) ToListByGroupId(groupId int, pageSize int, pageIndex int) collections.PageList[vo.TaskDO] {
	page := repository.task.Where("TaskGroupId = ?", groupId).Select("Id", "Caption", "Progress", "Status", "StartAt", "CreateAt", "ClientIp", "RunSpeed", "RunAt").Desc("CreateAt").ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[vo.TaskDO](page.List, page.RecordCount)
}

func (repository taskGroupRepository) ToListByClientId(clientId int64) collections.List[taskGroup.DomainObject] {
	lst := repository.ToList()
	return lst.Where(func(item taskGroup.DomainObject) bool {
		return item.Task.Client.Id == clientId && item.Task.StartAt.UnixMicro() < time.Now().UnixMicro()
	}).ToList()
}

func (repository taskGroupRepository) GetTaskGroupCount() int64 {
	return repository.taskGroup.Count()
}

func (repository taskGroupRepository) ToFinishList(taskGroupId int, top int) collections.List[vo.TaskDO] {
	lstPO := repository.task.Where("TaskGroupId = ? and (Status = ? or Status = ?)", taskGroupId, eumTaskType.Success, eumTaskType.Fail).Desc("CreateAt").Limit(top).ToList()
	var lstDO collections.List[vo.TaskDO]
	lstPO.MapToList(&lstDO)
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
	for _, do := range lst.ToArray() {
		po := mapper.Single[model.TaskGroupPO](do)
		repository.taskGroup.Where("Id = ?", do.Id).Update(po)
	}
}

func (repository taskGroupRepository) GetCanSchedulerTaskGroup(jobsName []string, ts time.Duration, count int, client vo.ClientVO) collections.List[vo.TaskDO] {
	getLocker := repository.redis.Lock.GetLocker("FSS_Scheduler", 5*time.Second)
	if !getLocker.TryLock() {
		exception.ThrowRefuseException("加锁失败")
	}
	defer getLocker.ReleaseLock()
	lstSchedulerTaskGroup := repository.ToList().Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable &&
			collections.NewList(jobsName...).Contains(item.JobName) &&
			item.Task.Status == eumTaskType.None &&
			item.Task.StartAt.UnixMicro() < time.Now().Add(ts).UnixMicro() &&
			item.Task.Client.Id == 0
	}).OrderBy(func(item taskGroup.DomainObject) any {
		return item.StartAt.UnixMicro()
	}).Take(count)

	var lst collections.List[vo.TaskDO]
	for _, taskGroupDO := range lstSchedulerTaskGroup.ToArray() {
		// 设为调度状态
		taskGroupDO.Scheduler(client)
		repository.Save(taskGroupDO)
		// 如果不相等，说明被其它客户端拿了
		if taskGroupDO.Task.Client.Id == client.Id {
			lst.Add(taskGroupDO.Task)
		}
	}
	return lst
}

func (repository taskGroupRepository) ToUnRunCount() int {
	return repository.ToList().Where(func(item taskGroup.DomainObject) bool {
		return item.Task.Status == eumTaskType.None || item.Task.Status == eumTaskType.Scheduler || item.Task.CreateAt.UnixMicro() < time.Now().UnixMicro()
	}).Count()
}

func (repository taskGroupRepository) ToSchedulerWorkingList() collections.List[taskGroup.DomainObject] {
	return repository.ToList().Where(func(item taskGroup.DomainObject) bool {
		return item.Task.Status == eumTaskType.Scheduler || item.Task.Status == eumTaskType.Working
	}).ToList()
}

func (repository taskGroupRepository) ToFinishPageList(pageSize int, pageIndex int) collections.PageList[vo.TaskDO] {
	pageList := repository.task.Where("(Status = ? or Status = ?) and (CreateAt >= ?)", eumTaskType.Fail, eumTaskType.Success, time.Now().Add(-24*time.Hour)).
		Select("Id", "Caption", "Progress", "Status", "StartAt", "CreateAt", "ClientIp", "RunSpeed", "RunAt", "JobName").
		Desc("RunAt").ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[vo.TaskDO](pageList.List, pageList.RecordCount)
}

func (repository taskGroupRepository) GetTaskUnFinishList(jobsName []string, top int) collections.List[taskGroup.DomainObject] {
	return repository.ToList().Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable && collections.NewList(jobsName...).Contains(item.JobName) && item.Task.Status != eumTaskType.Success && item.Task.Status != eumTaskType.Fail
	}).OrderBy(func(item taskGroup.DomainObject) any {
		return item.StartAt.UnixMicro()
	}).Take(top).ToList()
}

func (repository taskGroupRepository) GetEnableTaskList(status eumTaskType.Enum, pageSize int, pageIndex int) collections.PageList[vo.TaskDO] {
	lstTaskGroup := repository.ToList().Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable
	}).ToList()

	if status != eumTaskType.None {
		lstTaskGroup = lstTaskGroup.Where(func(item taskGroup.DomainObject) bool {
			return item.Task.Status == status
		}).ToList()
	}

	lstTaskGroup = lstTaskGroup.OrderBy(func(item taskGroup.DomainObject) any {
		return item.JobName
	}).ToList()

	var lst collections.List[vo.TaskDO]
	lstTaskGroup.MapToList(&lst)
	return lst.ToPageList(pageSize, pageIndex)
}
