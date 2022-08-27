package interfaces

import (
	"fss/application/clients/clientApp"
	"fss/application/tasks/taskGroupApp"
	"fss/application/tasks/taskGroupApp/request"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
)

type TestController struct {
}

// Pull 测试自研的API框架示例
func (r TestController) Pull(dto request.PullDTO) core.ApiResponse[collections.List[request.TaskDTO]] {

	lstTaskDTO := taskGroupApp.Pull(r.GetClient(), dto)
	apiResponse := core.Success("成功", lstTaskDTO)

	return apiResponse
}

func (r TestController) GetClient() clientApp.DTO {
	return clientApp.DTO{}
}
