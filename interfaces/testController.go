package interfaces

type TestController struct {
}

//// Pull 测试自研的API框架示例
//func (r TestController) Pull(dto request.PullDTO) core.ApiResponse[collections.List[request.TaskDTO]] {
//
//	lstTaskDTO := taskGroupApp.Pull(r.GetClient(), dto)
//	apiResponse := core.Success("成功", lstTaskDTO)
//
//	return apiResponse
//}
//
//func (r TestController) GetClient() clientApp.DTO {
//	return clientApp.DTO{}
//}
