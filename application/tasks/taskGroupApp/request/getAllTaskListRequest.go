package request

import "fss/domain/_/eumTaskType"

type GetAllTaskListRequest struct {
	Status    eumTaskType.Enum // 状态
	PageSize  int
	PageIndex int
}
