package repository

import (
	"fss/domain/log/taskLog"
	"fss/infrastructure/repository/model"
	"github.com/farseer-go/collections"
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
	container.Use[taskLog.Repository](func() taskLog.Repository {
		var repository taskLogRepository
		elasticSearch.InitContext(&repository, "es")
		return repository
	}).Transient().Register()
}

type taskLogRepository struct {
	TaskLog elasticSearch.IndexSet[model.TaskLogPO] `es:"index=run_log_yyyy_MM;alias=run_log"`
}

func NewTaskLogRepository() taskLogRepository {
	return container.Resolve[taskLog.Repository]().(taskLogRepository)
}

func (repository taskLogRepository) GetList(jobName string, logLevel eumLogLevel.Enum, pageSize int, pageIndex int) collections.List[taskLog.DomainObject] {
	pageList := repository.TaskLog.Where("JobName", jobName).Where("LogLevel", strconv.Itoa(int(logLevel))).ToPageList(pageSize, pageIndex)
	return mapper.ToList[taskLog.DomainObject](pageList)
}

func (repository taskLogRepository) Add(taskLogDO taskLog.DomainObject) {
	po := mapper.Single[model.TaskLogPO](&taskLogDO)
	queue.Push("TaskLogQueue", po)
}

func (repository taskLogRepository) AddBatch(lstPO collections.List[model.TaskLogPO]) {
	err := repository.TaskLog.InsertList(lstPO)
	if err != nil {
		exception.ThrowRefuseException("批量添加报错")
	}
}
