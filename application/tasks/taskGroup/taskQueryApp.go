package taskGroup

import (
	"fss/application/tasks/taskGroup/request"
	"fss/domain/clients/client"
	"fss/domain/tasks/taskGroup"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/mapper"
)

type taskQueryApp struct {
	repository       taskGroup.Repository
	clientRepository client.Repository
}

func NewTaskQueryApp() *taskQueryApp {
	return &taskQueryApp{
		repository:       container.Resolve[taskGroup.Repository](),
		clientRepository: container.Resolve[client.Repository](),
	}
}

// ToEntity 获取任务组信息
func (r *taskQueryApp) ToEntity(request request.OnlyIdRequest) DTO {
	do := r.repository.ToEntity(request.Id)
	return mapper.Single[DTO](do)
}

// ToList 获取所有任务组中的任务
func (r *taskQueryApp) ToList() collections.List[taskGroup.DomainObject] {
	lstTaskGroup := r.repository.ToList()
	for i := 0; i < lstTaskGroup.Count(); i++ {
		do := lstTaskGroup.Index(i)
		if do.Task.Id < 1 {
			do.CreateTask()
			r.repository.Save(do)
		}
		lstTaskGroup.Set(i, do)
	}
	return lstTaskGroup
}

// GetTaskList ToList 获取指定任务组的任务列表（FOPS）
func (r *taskQueryApp) GetTaskList(getTaskListRequest request.GetTaskListRequest) collections.PageList[request.TaskDTO] {
	page := r.repository.ToListByGroupId(getTaskListRequest.GroupId, getTaskListRequest.PageSize, getTaskListRequest.PageIndex)
	return mapper.ToPageList[request.TaskDTO](page.List, page.RecordCount)
}

// TodayFailCount 今日执行失败数量
func (r *taskQueryApp) TodayFailCount() int64 {
	return r.repository.TodayFailCount()
}

// GetTaskGroupCount 获取任务组数量
func (r *taskQueryApp) GetTaskGroupCount() int64 {
	return r.repository.GetTaskGroupCount()
}

// ToUnRunCount 获取未执行的任务数量
func (r *taskQueryApp) ToUnRunCount() int {
	return r.repository.ToUnRunCount()
}

// ToFinishPageList 获取已完成的任务列表
func (r *taskQueryApp) ToFinishPageList(pageSizeRequest request.PageSizeRequest) collections.PageList[request.TaskDTO] {
	page := r.repository.ToFinishPageList(pageSizeRequest.PageSize, pageSizeRequest.PageIndex)
	return mapper.ToPageList[request.TaskDTO](page.List, page.RecordCount)
}

// GetTaskUnFinishList 获取进行中的任务
func (r *taskQueryApp) GetTaskUnFinishList(onlyTopRequest request.OnlyTopRequest) collections.List[request.TaskDTO] {
	lstClient := r.clientRepository.ToList()
	if lstClient.IsEmpty() {
		return collections.NewList[request.TaskDTO]()
	}

	var lstJobs []string
	lstClient.SelectMany(&lstJobs, func(item client.DomainObject) any {
		return item.Jobs
	})

	taskUnFinishList := r.repository.GetTaskUnFinishList(lstJobs, onlyTopRequest.Top)

	var lstTask collections.List[request.TaskDTO]
	taskUnFinishList.Select(&lstTask, func(item taskGroup.DomainObject) any {
		return mapper.Single[request.TaskDTO](item.Task)
	})

	return lstTask
}

// GetEnableTaskList 获取在用的任务组
func (r *taskQueryApp) GetEnableTaskList(getAllTaskListRequest request.GetAllTaskListRequest) collections.PageList[request.TaskDTO] {
	var lst = r.repository.GetEnableTaskList(getAllTaskListRequest.Status, getAllTaskListRequest.PageSize, getAllTaskListRequest.PageIndex)
	return mapper.ToPageList[request.TaskDTO](lst.List, lst.RecordCount)
}
