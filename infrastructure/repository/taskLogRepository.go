package repository

import (
	"fss/domain/log/taskLog"
	"fss/infrastructure/repository/model"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/elasticSearch"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/queue"
	"strconv"
)

func RegisterTaskLogRepository() {
	// 注册仓储
	container.Register(func() taskLog.Repository {
		repository := data.NewContext[taskLogRepository]("default")
		repository.taskLogES = elasticSearch.NewContext[elasticSearchContext]("default").taskLog
		return repository

	})

}

type taskLogRepository struct {
	taskLog   data.TableSet[model.TaskLogPO] `data:"name=run_log"`
	taskLogES elasticSearch.IndexSet[model.TaskLogPO]
}
type elasticSearchContext struct {
	taskLog elasticSearch.IndexSet[model.TaskLogPO] `es:"name=run_log"`
}

func NewTaskLogRepository() taskLogRepository {
	return container.Resolve[taskLog.Repository]().(taskLogRepository)
}

func (repository taskLogRepository) GetList(jobName string, logLevel eumLogLevel.Enum, pageSize int, pageIndex int) collections.List[taskLog.DomainObject] {
	pageList := repository.taskLogES.Where("JobName", jobName).Where("LogLevel", strconv.Itoa(int(logLevel))).ToPageList(pageSize, pageIndex)
	return mapper.ToList[taskLog.DomainObject](pageList)
}

func (repository taskLogRepository) Add(taskLogDO taskLog.DomainObject) {
	po := mapper.Single[model.TaskLogPO](taskLogDO)
	queue.Push("TaskLogQueue", po)
}

func (repository taskLogRepository) AddBatch(lstPO collections.List[model.TaskLogPO]) {
	isSuccess, _ := repository.taskLogES.InsertList(lstPO)
	if !isSuccess {
		exception.ThrowRefuseException("批量添加报错")
	}
}
