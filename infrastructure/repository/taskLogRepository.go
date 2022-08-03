package repository

import (
	"fss/domain/log/taskLog"
	"fss/infrastructure/repository/model"
	"github.com/farseernet/farseer.go/core/container"
	"github.com/farseernet/farseer.go/core/eumLogLevel"
	"github.com/farseernet/farseer.go/mapper"
	"github.com/farseernet/farseer.go/mq/queue"
)

func init() {
	// 注册仓储
	_ = container.Register(func() taskLog.Repository {
		return &taskLogRepository{
			//data.Init[context.MysqlContext]("default").Admin,
		}
	})
}

type taskLogRepository struct {
	//data.TableSet[model.TaskLogPO]
}

func NewTaskLogRepository() *taskLogRepository {
	return &taskLogRepository{}
}

func (repository taskLogRepository) GetList(jobName string, logLevel eumLogLevel.Enum, pageSize int, pageIndex int) []taskLog.DomainObject {
	//TODO implement me
	panic("implement me")
}

func (repository taskLogRepository) Add(taskLogDO taskLog.DomainObject) {
	po := mapper.Single[model.TaskLogPO](taskLogDO)
	queue.Push("TaskLogQueue", po)
}

func (repository taskLogRepository) AddBatch(lstPO []model.TaskLogPO) {
	// todo
	panic("AddBatch未实现")
}
