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
		Id:         parse.Convert(r.HttpContext.Header.GetValue("ClientId"), int64(0)),
		Ip:         strings.Split(r.HttpContext.Header.GetValue("ClientIp"), ",")[0],
		Name:       r.HttpContext.Header.GetValue("ClientName"),
		Jobs:       strings.Split(r.HttpContext.Header.GetValue("ClientJobs"), ","),
		ActivateAt: time.Now(),
	}
	clientApp.UpdateClient(dto)
	return dto
}

//func (r *TaskController) JobInvoke() {
//	// 读取结构数据
//	var dto request.JobInvokeRequest
//	_ = r.BindJSON(&dto)
//
//	result := taskGroupApp.JobInvoke(r.GetClient(), dto)
//	apiResponse := core.ApiResponseStringSuccess(result)
//
//	// 响应数据
//	r.Data["json"] = &apiResponse
//	_ = r.ServeJSON()
//}
//
//func (r *TaskController) Pull() {
//	// 读取结构数据
//	var dto request.PullRequest
//	_ = r.BindJSON(&dto)
//
//	lstTaskDTO := taskGroupApp.Pull(r.GetClient(), dto)
//	apiResponse := core.Success("成功", lstTaskDTO)
//
//	// 响应数据
//	r.Data["json"] = &apiResponse
//	_ = r.ServeJSON()
//}
