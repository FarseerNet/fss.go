package request

import "fss/domain/_/eumTaskType"

type GetAllTaskListRequest struct {
	// 状态
	Status    eumTaskType.Enum
	PageSize  int
	PageIndex int
}
