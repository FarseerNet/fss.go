package repository

import (
	"fss/domain/_/eumTaskType"
	"fss/domain/tasks/taskGroup"
	"fss/domain/tasks/taskGroup/vo"
	"fss/infrastructure/repository/model"
	"github.com/farseer-go/cache"
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

func RegisterTaskGroupRepository() {
	// 注册仓储
	container.Register(func() taskGroup.Repository {
		repository := data.NewContext[taskGroupRepository]("default")
		repository.redis = redis.NewClient("default")

		// 多级缓存
		repository.cacheManage = cache.GetCacheManage[taskGroup.DomainObject]("FSS_TaskGroup")
		repository.cacheManage.EnableItemNullToLoadALl()
		repository.cacheManage.SetSource(func() collections.List[taskGroup.DomainObject] {
			var lst collections.List[taskGroup.DomainObject]
			repository.taskGroup.ToList().MapToList(&lst)
			return lst
		})
		return repository
	})
}

type taskGroupRepository struct {
	taskGroup   data.TableSet[model.TaskGroupPO] `data:"name=task_group"`
	task        data.TableSet[model.TaskPO]      `data:"name=task"`
	redis       *redis.Client
	cacheManage cache.CacheManage[taskGroup.DomainObject]
}

func (repository taskGroupRepository) ToList() collections.List[taskGroup.DomainObject] {
	return repository.cacheManage.Get()
}

func (repository taskGroupRepository) ToEntity(taskGroupId int) taskGroup.DomainObject {
	item, _ := repository.cacheManage.GetItem(taskGroupId)
	return item
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

func (repository taskGroupRepository) ToListByGroupId(groupId int, pageSize int, pageIndex int) collections.PageList[vo.TaskEO] {
	page := repository.task.Where("TaskGroupId = ?", groupId).Select("Id", "Caption", "Progress", "Status", "StartAt", "CreateAt", "ClientIp", "RunSpeed", "RunAt").Desc("CreateAt").ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[vo.TaskEO](page.List, page.RecordCount)
}

func (repository taskGroupRepository) ToListByClientId(clientId int64) collections.List[taskGroup.DomainObject] {
	lst := repository.ToList()
	return lst.Where(func(item taskGroup.DomainObject) bool {
		return item.Task.Client.Id == clientId && item.Task.StartAt.UnixMicro() < time.Now().UnixMicro()
	}).ToList()
}

func (repository taskGroupRepository) GetTaskGroupCount() int64 {
	return int64(repository.cacheManage.Count())
}

func (repository taskGroupRepository) ToFinishList(taskGroupId int, top int) collections.List[vo.TaskEO] {
	lstPO := repository.task.Where("TaskGroupId = ? and (Status = ? or Status = ?)", taskGroupId, eumTaskType.Success, eumTaskType.Fail).Desc("CreateAt").Limit(top).ToList()
	var lstDO collections.List[vo.TaskEO]
	lstPO.MapToList(&lstDO)
	return lstDO
}

func (repository taskGroupRepository) AddTask(taskDO vo.TaskEO) {
	po := mapper.Single[model.TaskPO](taskDO)
	repository.task.Insert(&po)
}

func (repository taskGroupRepository) Add(do taskGroup.DomainObject) taskGroup.DomainObject {
	po := mapper.Single[model.TaskGroupPO](do)
	repository.taskGroup.Insert(&po)
	do.Id = po.Id
	repository.cacheManage.SaveItem(do)
	return do
}

func (repository taskGroupRepository) Save(do taskGroup.DomainObject) {
	repository.cacheManage.SaveItem(do)
}

func (repository taskGroupRepository) Delete(taskGroupId int) {
	repository.task.Where("TaskGroupId = ?", taskGroupId).Delete()
	repository.taskGroup.Where("Id = ?", taskGroupId).Delete()
	repository.cacheManage.Remove(strconv.Itoa(taskGroupId))
}

func (repository taskGroupRepository) SyncToData() {
	lst := repository.ToList()
	for _, do := range lst.ToArray() {
		po := mapper.Single[model.TaskGroupPO](do)
		repository.taskGroup.Where("Id = ?", do.Id).Update(po)
	}
}

func (repository taskGroupRepository) GetCanSchedulerTaskGroup(jobsName []string, ts time.Duration, count int, client vo.ClientVO) collections.List[vo.TaskEO] {
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

	var lst collections.List[vo.TaskEO]
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

func (repository taskGroupRepository) ToFinishPageList(pageSize int, pageIndex int) collections.PageList[vo.TaskEO] {
	pageList := repository.task.Where("(Status = ? or Status = ?) and (CreateAt >= ?)", eumTaskType.Fail, eumTaskType.Success, time.Now().Add(-24*time.Hour)).
		Select("Id", "Caption", "Progress", "Status", "StartAt", "CreateAt", "ClientIp", "RunSpeed", "RunAt", "JobName").
		Desc("RunAt").ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[vo.TaskEO](pageList.List, pageList.RecordCount)
}

func (repository taskGroupRepository) GetTaskUnFinishList(jobsName []string, top int) collections.List[taskGroup.DomainObject] {
	return repository.ToList().Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable && collections.NewList(jobsName...).Contains(item.JobName) && item.Task.Status != eumTaskType.Success && item.Task.Status != eumTaskType.Fail
	}).OrderBy(func(item taskGroup.DomainObject) any {
		return item.StartAt.UnixMicro()
	}).Take(top).ToList()
}

func (repository taskGroupRepository) GetEnableTaskList(status eumTaskType.Enum, pageSize int, pageIndex int) collections.PageList[vo.TaskEO] {
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

	var lst collections.List[vo.TaskEO]
	lstTaskGroup.MapToList(&lst)
	return lst.ToPageList(pageSize, pageIndex)
}
