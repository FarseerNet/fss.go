package repository

import (
	"fss/domain/log/taskLog"
	"fss/infrastructure/repository/model"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/queue"
)

func RegisterTaskLogRepository() {
	// 注册仓储
	container.Register(func() taskLog.Repository {
		return data.NewContext[taskLogRepository]("default")
	})
}

type taskLogRepository struct {
	data.TableSet[model.TaskLogPO] `data:"name=run_log"`
}

func NewTaskLogRepository() taskLogRepository {
	return container.Resolve[taskLog.Repository]().(taskLogRepository)
}

func (repository taskLogRepository) GetList(jobName string, logLevel eumLogLevel.Enum, pageSize int, pageIndex int) collections.List[taskLog.DomainObject] {
	//TODO implement me
	panic("ES日志查询未实现")
}

func (repository taskLogRepository) Add(taskLogDO taskLog.DomainObject) {
	po := mapper.Single[model.TaskLogPO](taskLogDO)
	queue.Push("TaskLogQueue", po)
}

func (repository taskLogRepository) AddBatch(lstPO collections.List[model.TaskLogPO]) {
	// todo
	exception.ThrowRefuseException("AddBatch未实现")
}
