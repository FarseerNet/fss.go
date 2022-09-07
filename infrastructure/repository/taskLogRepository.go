package repository

import (
	"fss/domain/log/taskLog"
	"fss/infrastructure/repository/model"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/elasticSearch"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/queue"
	"strconv"
)

func RegisterTaskLogRepository() {
	// 注册仓储
	container.Register(func() taskLog.Repository {
		var repository taskLogRepository
		data.InitContext(&repository, "default")
		elasticSearch.InitContext(&repository, "es")
		return repository
	})
}

type taskLogRepository struct {
	TaskLog   data.TableSet[model.TaskLogPO]          `data:"name=run_log"`
	TaskLogES elasticSearch.IndexSet[model.TaskLogPO] `es:"name=run_log"`
}

func NewTaskLogRepository() taskLogRepository {
	return container.Resolve[taskLog.Repository]().(taskLogRepository)
}

func (repository taskLogRepository) GetList(jobName string, logLevel eumLogLevel.Enum, pageSize int, pageIndex int) collections.List[taskLog.DomainObject] {
	pageList := repository.TaskLogES.Where("JobName", jobName).Where("LogLevel", strconv.Itoa(int(logLevel))).ToPageList(pageSize, pageIndex)
	return mapper.ToList[taskLog.DomainObject](pageList)
}

func (repository taskLogRepository) Add(taskLogDO taskLog.DomainObject) {
	po := mapper.Single[model.TaskLogPO](&taskLogDO)
	queue.Push("TaskLogQueue", po)
}

func (repository taskLogRepository) AddBatch(lstPO collections.List[model.TaskLogPO]) {
	err := repository.TaskLogES.InsertList(lstPO)
	if err != nil {
		exception.ThrowRefuseException("批量添加报错")
	}
}
