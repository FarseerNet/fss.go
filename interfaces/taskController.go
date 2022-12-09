package interfaces

import (
	"fss/application/clients/clientApp"
	"fss/application/tasks/taskGroupApp"
	"fss/application/tasks/taskGroupApp/request"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/webapi/controller"
	"strings"
	"time"
)

type TaskController struct {
	controller.BaseController
}

func (r *TaskController) JobInvoke(dto request.JobInvokeRequest) string {
	return taskGroupApp.JobInvoke(r.GetClient(), dto)
}

func (r *TaskController) Pull(dto request.PullRequest) collections.List[request.TaskDTO] {
	return taskGroupApp.Pull(r.GetClient(), dto)
}

// GetClient 获取头部信息，并更新客户端
func (r *TaskController) GetClient() clientApp.DTO {
	dto := clientApp.DTO{
		Id:         parse.Convert(r.HttpContext.Header.GetValue("Clientid"), int64(0)),
		Ip:         strings.Split(r.HttpContext.Header.GetValue("Clientip"), ",")[0],
		Name:       r.HttpContext.Header.GetValue("Clientname"),
		Jobs:       strings.Split(r.HttpContext.Header.GetValue("Clientjobs"), ","),
		ActivateAt: time.Now(),
	}
	clientApp.UpdateClient(dto)
	return dto
}
