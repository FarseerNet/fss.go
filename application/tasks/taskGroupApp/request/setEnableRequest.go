package request

// SetEnableRequest 设置任务组状态
type SetEnableRequest struct {
	Id       int  // 任务组ID
	IsEnable bool // 任务组状态
}
