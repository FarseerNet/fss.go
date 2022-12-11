package interfaces

import (
	"fss/application/clients/clientApp"
	"fss/application/tasks/taskGroupApp"
	"fss/application/tasks/taskGroupApp/request"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/webapi/controller"
	"time"
)

type TaskController struct {
	controller.BaseController
	Client clientApp.DTO `webapi:"header"`
}

func (r *TaskController) JobInvoke(dto request.JobInvokeRequest) string {
	return taskGroupApp.JobInvoke(r.Client, dto)
}

func (r *TaskController) Pull(dto request.PullRequest) collections.List[request.TaskDTO] {
	return taskGroupApp.Pull(r.Client, dto)
}

func (r *TaskController) OnActionExecuting() {
	r.Client.ActivateAt = time.Now()
	clientApp.UpdateClient(r.Client)
}
func (r *TaskController) OnActionExecuted() {
}
