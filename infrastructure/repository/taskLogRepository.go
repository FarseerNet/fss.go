package repository

import (
	"fss/domain/log/taskLog"
	"fss/infrastructure/repository/model"
	"github.com/farseernet/farseer.go/core/container"
	"github.com/farseernet/farseer.go/core/eumLogLevel"
	"github.com/farseernet/farseer.go/data"
	"github.com/farseernet/farseer.go/exception"
	"github.com/farseernet/farseer.go/mapper"
	"github.com/farseernet/farseer.go/mq/queue"
)

func RegisterTaskLogRepository() {
	// 注册仓储
	_ = container.Register(func() taskLog.Repository {
		return data.NewContext[taskLogRepository]("default")
	})
}

type taskLogRepository struct {
	data.TableSet[model.TaskLogPO] `data:"name=run_log"`
}

func NewTaskLogRepository() taskLogRepository {
	return container.Resolve[taskLog.Repository]().(taskLogRepository)
}

func (repository taskLogRepository) GetList(jobName string, logLevel eumLogLevel.Enum, pageSize int, pageIndex int) []taskLog.DomainObject {
	//TODO implement me
	panic("ES日志查询未实现")
}

func (repository taskLogRepository) Add(taskLogDO taskLog.DomainObject) {
	po := mapper.Single[model.TaskLogPO](taskLogDO)
	queue.Push("TaskLogQueue", po)
}

func (repository taskLogRepository) AddBatch(lstPO []model.TaskLogPO) {
	// todo
	exception.ThrowRefuseException("AddBatch未实现")
}
