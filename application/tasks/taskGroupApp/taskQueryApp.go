package taskGroupApp

import (
	"fss/application/tasks/taskGroupApp/request"
	"fss/domain/clients/client"
	"fss/domain/tasks/taskGroup"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/mapper"
)

// ToEntity 获取任务组信息
func ToEntity(request request.OnlyIdRequest) DTO {
	repository := container.Resolve[taskGroup.Repository]()
	do := repository.ToEntity(request.Id)
	return mapper.Single[DTO](do)
}

// ToList 获取所有任务组中的任务
func ToList() collections.List[taskGroup.DomainObject] {
	repository := container.Resolve[taskGroup.Repository]()
	lstTaskGroup := repository.ToList()
	for i := 0; i < lstTaskGroup.Count(); i++ {
		do := lstTaskGroup.Index(i)
		if do.Task.Id < 1 {
			do.CreateTask()
			repository.Save(do)
		}
		lstTaskGroup.Set(i, do)
	}
	return lstTaskGroup
}

// GetTaskList ToList 获取指定任务组的任务列表（FOPS）
func GetTaskList(getTaskListRequest request.GetTaskListRequest) collections.PageList[request.TaskDTO] {
	repository := container.Resolve[taskGroup.Repository]()
	page := repository.ToListByGroupId(getTaskListRequest.GroupId, getTaskListRequest.PageSize, getTaskListRequest.PageIndex)
	return mapper.ToPageList[request.TaskDTO](page.List, page.RecordCount)
}

// TodayFailCount 今日执行失败数量
func TodayFailCount() int64 {
	repository := container.Resolve[taskGroup.Repository]()
	return repository.TodayFailCount()
}

// GetTaskGroupCount 获取任务组数量
func GetTaskGroupCount() int64 {
	repository := container.Resolve[taskGroup.Repository]()
	return repository.GetTaskGroupCount()
}

// ToUnRunCount 获取未执行的任务数量
func ToUnRunCount() int {
	repository := container.Resolve[taskGroup.Repository]()
	return repository.ToUnRunCount()
}

// ToFinishPageList 获取已完成的任务列表
func ToFinishPageList(pageSizeRequest request.PageSizeRequest) collections.PageList[request.TaskDTO] {
	repository := container.Resolve[taskGroup.Repository]()
	page := repository.ToFinishPageList(pageSizeRequest.PageSize, pageSizeRequest.PageIndex)
	return mapper.ToPageList[request.TaskDTO](page.List, page.RecordCount)
}

// GetTaskUnFinishList 获取进行中的任务
func GetTaskUnFinishList(onlyTopRequest request.OnlyTopRequest) collections.List[request.TaskDTO] {
	repository := container.Resolve[taskGroup.Repository]()
	clientRepository := container.Resolve[client.Repository]()
	lstClient := clientRepository.ToList()
	if lstClient.IsEmpty() {
		return collections.NewList[request.TaskDTO]()
	}

	var lstJobs []string
	lstClient.SelectMany(&lstJobs, func(item client.DomainObject) any {
		return item.Jobs
	})

	taskUnFinishList := repository.GetTaskUnFinishList(lstJobs, onlyTopRequest.Top)

	var lstTask collections.List[request.TaskDTO]
	taskUnFinishList.Select(&lstTask, func(item taskGroup.DomainObject) any {
		return mapper.Single[request.TaskDTO](item.Task)
	})

	return lstTask
}

// GetEnableTaskList 获取在用的任务组
func GetEnableTaskList(getAllTaskListRequest request.GetAllTaskListRequest) collections.PageList[request.TaskDTO] {
	repository := container.Resolve[taskGroup.Repository]()
	var lst = repository.GetEnableTaskList(getAllTaskListRequest.Status, getAllTaskListRequest.PageSize, getAllTaskListRequest.PageIndex)
	return mapper.ToPageList[request.TaskDTO](lst.List, lst.RecordCount)
}
