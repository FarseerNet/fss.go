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
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/redis"
	"strconv"
	"time"
)

type taskGroupRepository struct {
	TaskGroup   data.TableSet[model.TaskGroupPO] `data:"name=task_group"`
	Task        data.TableSet[model.TaskPO]      `data:"name=task"`
	redis       *redis.Client
	cacheManage cache.CacheManage[taskGroup.DomainObject]
}

func RegisterTaskGroupRepository() {
	// 注册仓储
	container.Register(func() taskGroup.Repository {
		repository := data.NewContext[taskGroupRepository]("default")
		repository.redis = redis.NewClient("default")

		// 多级缓存
		repository.cacheManage = cache.GetCacheManage[taskGroup.DomainObject]("FSS_TaskGroup")
		repository.cacheManage.SetListSource(func() collections.List[taskGroup.DomainObject] {
			var lst collections.List[taskGroup.DomainObject]
			repository.TaskGroup.ToList().MapToList(&lst)
			return lst
		})
		repository.cacheManage.SetItemSource(func(cacheId any) (taskGroup.DomainObject, bool) {
			po := repository.TaskGroup.Where("Id = ?", cacheId).ToEntity()
			if po.Id > 0 {
				return mapper.Single[taskGroup.DomainObject](&po), true
			}
			var do taskGroup.DomainObject
			return do, false
		})
		return *repository
	})
}

func (repository taskGroupRepository) ToList() collections.List[taskGroup.DomainObject] {
	return repository.cacheManage.Get()
}

func (repository taskGroupRepository) ToEntity(taskGroupId int) taskGroup.DomainObject {
	item, _ := repository.cacheManage.GetItem(taskGroupId)
	return item
}

func (repository taskGroupRepository) TodayFailCount() int64 {
	return repository.Task.Where("status = ? and create_at >= ?", eumTaskType.Fail, dateTime.Now().Date().ToTime()).Count()
}

func (repository taskGroupRepository) ToTaskSpeedList(taskGroupId int) []int64 {
	lstPO := repository.Task.Where("task_group_id = ? and status = ?", taskGroupId, eumTaskType.Success).Desc("create_at").Select("RunSpeed").Limit(100).ToList()
	var lstSpeed []int64
	lstPO.Select(&lstSpeed, func(item model.TaskPO) any {
		return item.RunSpeed
	})
	return lstSpeed
}

func (repository taskGroupRepository) ToListByGroupId(groupId int, pageSize int, pageIndex int) collections.PageList[vo.TaskEO] {
	page := repository.Task.Where("task_group_id = ?", groupId).Select("Id", "Caption", "Progress", "Status", "StartAt", "CreateAt", "ClientIp", "RunSpeed", "RunAt").Desc("create_at").ToPageList(pageSize, pageIndex)
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
	lstPO := repository.Task.Where("task_group_id = ? and (status = ? or status = ?)", taskGroupId, eumTaskType.Success, eumTaskType.Fail).Desc("create_at").Limit(top).ToList()
	lstDO := collections.NewList[vo.TaskEO]()
	for _, taskPO := range lstPO.ToArray() {
		task := mapper.Single[vo.TaskEO](&taskPO)
		lstDO.Add(task)
	}
	lstPO.MapToList(&lstDO)
	return lstDO
}

func (repository taskGroupRepository) AddTask(taskDO vo.TaskEO) {
	po := mapper.Single[model.TaskPO](&taskDO)
	po.ClientId = taskDO.Client.Id
	po.ClientName = taskDO.Client.Name
	po.ClientIp = taskDO.Client.Ip
	repository.Task.Insert(&po)
}

func (repository taskGroupRepository) Add(do *taskGroup.DomainObject) {
	po := mapper.Single[model.TaskGroupPO](do)
	repository.TaskGroup.Insert(&po)
	do.Id = po.Id
	repository.cacheManage.SaveItem(*do)
}

func (repository taskGroupRepository) Save(do taskGroup.DomainObject) {
	if do.Id == 0 {
		flog.Errorf("发现taskGroup.Id=0的数据，val=%v", do)
	}
	repository.cacheManage.SaveItem(do)
}

func (repository taskGroupRepository) Delete(taskGroupId int) {
	repository.Task.Where("task_group_id = ?", taskGroupId).Delete()
	repository.TaskGroup.Where("Id = ?", taskGroupId).Delete()
	repository.cacheManage.Remove(strconv.Itoa(taskGroupId))
}

func (repository taskGroupRepository) GetCanSchedulerTaskGroup(jobsName []string, ts time.Duration, count int, client vo.ClientVO) collections.List[vo.TaskEO] {
	getLocker := repository.redis.Lock.GetLocker("FSS_Scheduler", 5*time.Second)
	if !getLocker.TryLock() {
		flog.Warningf("调度任务时加锁失败，Job=%s，ClientIp=%s", collections.NewList(jobsName...).ToString(","), client.Ip)
		return collections.NewList[vo.TaskEO]()
	}
	defer getLocker.ReleaseLock()
	lstSchedulerTaskGroup := repository.ToList().Where(func(item taskGroup.DomainObject) bool {
		return item.CanScheduler(jobsName, ts)
	}).OrderBy(func(item taskGroup.DomainObject) any {
		return item.StartAt.UnixMicro()
	}).Take(count)

	lst := collections.NewList[vo.TaskEO]()
	for _, taskGroupDO := range lstSchedulerTaskGroup.ToArray() {
		// 设为调度状态
		taskGroupDO.Scheduler(client)
		repository.Save(taskGroupDO)
		// 如果不相等，说明被其它客户端拿了
		lst.Add(taskGroupDO.Task)
		if taskGroupDO.Task.TaskGroupId == 0 {
			flog.Errorf("发现taskGroupDO.Task.TaskGroupId=0的数据，val=%v", taskGroupDO.Task)
		}
		if taskGroupDO.Id == 0 {
			flog.Errorf("发现taskGroup.Id=0的数据，val=%v", taskGroupDO)
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
	pageList := repository.Task.Where("(status = ? or status = ?) and (create_at >= ?)", eumTaskType.Fail, eumTaskType.Success, time.Now().Add(-24*time.Hour)).
		Select("Id", "Caption", "Progress", "Status", "StartAt", "CreateAt", "ClientIp", "RunSpeed", "RunAt", "JobName").
		Desc("run_at").ToPageList(pageSize, pageIndex)
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
	for _, taskPO := range lstTaskGroup.ToArray() {
		task := mapper.Single[vo.TaskEO](&taskPO)
		lst.Add(task)
	}
	//lstTaskGroup.MapToList(&lst)
	return lst.ToPageList(pageSize, pageIndex)
}

// ClearFinish 清除成功的任务记录（1天前）
func (repository taskGroupRepository) ClearFinish(groupId int, taskId int) {
	repository.Task.Where("task_group_id = ? and (status = ? or status = ?) and create_at < ? and Id < ?", groupId, eumTaskType.Success, eumTaskType.Fail, time.Now().Add(-24*time.Hour), taskId).Delete()
}

// SaveToDb 保存到数据库
func (repository taskGroupRepository) SaveToDb(do taskGroup.DomainObject) {
	po := mapper.Single[model.TaskGroupPO](&do)
	if po.Id == 0 {
		flog.Errorf("发现taskGroup.Id=0的数据，val=%v", do)
	}

	repository.TaskGroup.Where("Id = ?", do.Id).Update(po)
}

// ToIdList 从数据库中读取数据
func (repository taskGroupRepository) ToIdList() []int {
	lst := repository.TaskGroup.Select("Id").ToList()
	var lstIds []int
	lst.Select(&lstIds, func(item model.TaskGroupPO) any {
		return item.Id
	})
	return lstIds
}
