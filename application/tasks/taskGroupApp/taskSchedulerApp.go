package taskGroupApp

import (
	"fss/application/clients/clientApp"
	"fss/application/tasks/taskGroupApp/request"
	"fss/domain/tasks/taskGroup"
	"fss/domain/tasks/taskGroup/vo"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/mapper"
	"time"
)

// Pull 任务调度
func Pull(clientDTO clientApp.DTO, dto request.PullDTO) collections.List[request.TaskDTO] {
	repository := container.Resolve[taskGroup.Repository]()
	clientVO := mapper.Single[vo.ClientVO](&clientDTO)
	if dto.TaskCount == 0 {
		dto.TaskCount = 3
	}

	s := 15 * time.Second
	lstDO := repository.GetCanSchedulerTaskGroup(clientDTO.Jobs, s, dto.TaskCount, clientVO)
	var lst collections.List[request.TaskDTO]
	lstDO.MapToList(&lst)

	if lst.Where(func(item request.TaskDTO) bool { return item.TaskGroupId == 0 }).Any() {
		flog.Errorf("发现taskGroup.Id=0的数据，count=%d", lstDO.Count())
	}
	return lst
}
