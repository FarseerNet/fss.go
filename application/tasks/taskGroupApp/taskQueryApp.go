package taskGroupApp

import (
	"fss/application/tasks/taskGroupApp/request"
	"fss/domain/clients/client"
	"fss/domain/tasks/taskGroup"
	"fss/domain/tasks/taskGroup/vo"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/mapper"
)

// ToEntity 获取任务组信息
func ToEntity(request request.OnlyIdRequest) DTO {
	repository := container.Resolve[taskGroup.Repository]()
	do := repository.ToEntity(request.Id)
	return mapper.Single[DTO](&do)
}

// ToList 获取所有任务组中的任务
func ToList() collections.List[taskGroup.DomainObject] {
	repository := container.Resolve[taskGroup.Repository]()
	return repository.ToList()
}

// TodayTaskFailCount 今日执行失败数量
func TodayTaskFailCount() int64 {
	repository := container.Resolve[taskGroup.Repository]()
	return repository.TodayFailCount()
}

// GetTaskGroupCount 获取任务组数量
func GetTaskGroupCount() int64 {
	repository := container.Resolve[taskGroup.Repository]()
	return repository.GetTaskGroupCount()
}

// GetTaskGroupUnRunCount 获取未执行的任务数量
func GetTaskGroupUnRunCount() int {
	repository := container.Resolve[taskGroup.Repository]()
	return repository.ToUnRunCount()
}

// GetTaskList ToList 获取指定任务组的任务列表（FOPS）
func GetTaskList(getTaskListRequest request.GetTaskListRequest) collections.PageList[vo.TaskEO] {
	repository := container.Resolve[taskGroup.Repository]()
	return repository.ToListByGroupId(getTaskListRequest.GroupId, getTaskListRequest.PageSize, getTaskListRequest.PageIndex)

	//return mapper.ToPageList[request.TaskDTO](page.List, page.RecordCount)
	//lst := collections.NewList[request.TaskDTO]()
	//for i := 0; i < page.List.Count(); i++ {
	//	item := page.List.Index(i)
	//	dto := mapper.Single[request.TaskDTO](&item)
	//	dto.ClientId = item.Client.Id
	//	dto.ClientIp = item.Client.Ip
	//	dto.ClientName = item.Client.Name
	//	lst.Add(dto)
	//}
	//return collections.NewPageList[request.TaskDTO](lst, page.RecordCount)
}

// GetTaskFinishList 获取已完成的任务列表
func GetTaskFinishList(pageSizeRequest request.PageSizeRequest) collections.PageList[vo.TaskEO] {
	repository := container.Resolve[taskGroup.Repository]()
	return repository.ToFinishPageList(pageSizeRequest.PageSize, pageSizeRequest.PageIndex)

	//return mapper.ToPageList[request.TaskDTO](page.List, page.RecordCount)
	//lst := collections.NewList[request.TaskDTO]()
	//for i := 0; i < page.List.Count(); i++ {
	//	item := page.List.Index(i)
	//	dto := mapper.Single[request.TaskDTO](&item)
	//	dto.ClientId = item.Client.Id
	//	dto.ClientIp = item.Client.Ip
	//	dto.ClientName = item.Client.Name
	//	lst.Add(dto)
	//}
	//return collections.NewPageList[request.TaskDTO](lst, page.RecordCount)
}

// GetTaskUnFinishList 获取进行中的任务
func GetTaskUnFinishList(onlyTopRequest request.OnlyTopRequest) collections.List[vo.TaskEO] {
	repository := container.Resolve[taskGroup.Repository]()
	clientRepository := container.Resolve[client.Repository]()
	lstClient := clientRepository.ToList()
	if lstClient.IsEmpty() {
		return collections.NewList[vo.TaskEO]()
	}

	var lstJobs []string
	lstClient.SelectMany(&lstJobs, func(item client.DomainObject) any {
		return item.Jobs
	})

	taskUnFinishList := repository.GetTaskUnFinishList(lstJobs, onlyTopRequest.Top)

	var lstTask collections.List[vo.TaskEO]
	taskUnFinishList.Select(&lstTask, func(item taskGroup.DomainObject) any {
		return item.Task
	})
	return lstTask
}

// GetEnableTaskList 获取在用的任务组
func GetEnableTaskList(getAllTaskListRequest request.GetAllTaskListRequest) collections.PageList[vo.TaskEO] {
	repository := container.Resolve[taskGroup.Repository]()
	return repository.GetEnableTaskList(getAllTaskListRequest.Status, getAllTaskListRequest.PageSize, getAllTaskListRequest.PageIndex)
	//return mapper.ToPageList[request.TaskDTO](page.List, page.RecordCount)

	//lst := collections.NewList[request.TaskDTO]()
	//for i := 0; i < page.List.Count(); i++ {
	//	item := page.List.Index(i)
	//	dto := mapper.Single[request.TaskDTO](&item)
	//	dto.ClientId = item.Client.Id
	//	dto.ClientIp = item.Client.Ip
	//	dto.ClientName = item.Client.Name
	//	lst.Add(dto)
	//}
	//return collections.NewPageList[request.TaskDTO](lst, page.RecordCount)
}
