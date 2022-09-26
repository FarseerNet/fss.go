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
	TaskLog elasticSearch.IndexSet[model.TaskLogPO] `es:"index=fsslog_yyyy_MM;alias=fsslog"`
}

func NewTaskLogRepository() taskLogRepository {
	return container.Resolve[taskLog.Repository]().(taskLogRepository)
}

func (repository taskLogRepository) GetList(jobName string, logLevel eumLogLevel.Enum, pageSize int, pageIndex int) collections.PageList[taskLog.DomainObject] {
	pageList := repository.TaskLog.Where("JobName", jobName).Where("LogLevel", logLevel).ToPageList(pageSize, pageIndex)
	var pageListDO collections.PageList[taskLog.DomainObject]
	pageList.MapToPageList(&pageListDO)
	return pageListDO
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
