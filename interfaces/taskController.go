package interfaces

import (
	"fss/application/clients/clientApp"
	"fss/application/tasks/taskGroupApp"
	"fss/application/tasks/taskGroupApp/request"
	"github.com/beego/beego/v2/server/web"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/stopwatch"
	"strings"
	"time"
)

type TaskController struct {
	web.Controller
}

func (r *TaskController) JobInvoke() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "JobInvoke，耗时：%s", sw.GetMillisecondsText())
	}()

	// 读取结构数据
	var dto request.JobInvokeDTO
	_ = r.BindJSON(&dto)

	result := taskGroupApp.JobInvoke(r.GetClient(), dto)
	apiResponse := core.ApiResponseStringSuccess(result)

	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

func (r *TaskController) Pull() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "Pull，耗时：%s", sw.GetMillisecondsText())
	}()

	// 读取结构数据
	var dto request.PullDTO
	_ = r.BindJSON(&dto)

	lstTaskDTO := taskGroupApp.Pull(r.GetClient(), dto)
	apiResponse := core.Success("成功", lstTaskDTO)

	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// GetClient 获取头部信息，并更新客户端
func (r *TaskController) GetClient() clientApp.DTO {
	dto := clientApp.DTO{
		Id:         parse.Convert(r.Ctx.Request.Header.Get("ClientId"), int64(0)),
		Ip:         strings.Split(r.Ctx.Request.Header.Get("ClientIp"), ",")[0],
		Name:       r.Ctx.Request.Header.Get("ClientName"),
		Jobs:       strings.Split(r.Ctx.Request.Header.Get("ClientJobs"), ","),
		ActivateAt: time.Now(),
	}
	clientApp.UpdateClient(dto)
	return dto
}
