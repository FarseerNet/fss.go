package taskGroup

import (
	"fss/application/clients/client"
	"fss/application/tasks/taskGroup/request"
	"fss/domain/tasks/taskGroup"
	"fss/domain/tasks/taskGroup/vo"
	"github.com/farseernet/farseer.go/core/container"
	"github.com/farseernet/farseer.go/mapper"
	"time"
)

type taskSchedulerApp struct {
	repository taskGroup.Repository
}

func NewTaskSchedulerApp() *taskSchedulerApp {
	return &taskSchedulerApp{repository: container.Resolve[taskGroup.Repository]()}
}

// Pull 任务调度
func (r *taskSchedulerApp) Pull(dto request.PullDTO) []request.TaskDTO {
	clientDTO := client.NewApp().GetClient()
	clientVO := mapper.Single[vo.ClientVO](clientDTO)
	if dto.TaskCount == 0 {
		dto.TaskCount = 3
	}

	s := 15 * time.Second
	lstDO := r.repository.GetCanSchedulerTaskGroup(clientDTO.Jobs, s, dto.TaskCount, clientVO)
	return mapper.Array[request.TaskDTO](lstDO)
}
