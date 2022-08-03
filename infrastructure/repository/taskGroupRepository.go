package repository

import (
	"fss/domain/_/eumTaskType"
	"fss/domain/tasks/taskGroup"
	"fss/domain/tasks/taskGroup/vo"
	"fss/infrastructure/repository/context"
	"fss/infrastructure/repository/model"
	"github.com/farseernet/farseer.go/core"
	"github.com/farseernet/farseer.go/core/container"
	"github.com/farseernet/farseer.go/data"
	"time"
)

func init() {
	// 注册仓储
	_ = container.Register(func() taskGroup.Repository {
		return &taskGroupRepository{
			data.Init[context.MysqlContext]("default").Admin,
		}
	})
}

type taskGroupRepository struct {
	data.TableSet[model.TaskGroupPO]
}

func (repository taskGroupRepository) ToEntity(taskGroupId int) taskGroup.DomainObject {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) TodayFailCount() int {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) ToTaskSpeedList(taskGroupId int) []int64 {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) ToList() []taskGroup.DomainObject {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) ToListByGroupId(groupId int, pageSize int, pageIndex int) []vo.TaskDO {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) ToListByClientId(clientId int64) []taskGroup.DomainObject {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) GetTaskGroupCount() int {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) ToFinishList(taskGroupId int, top int) []vo.TaskDO {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) AddTask(taskDO vo.TaskDO) {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) Add(do taskGroup.DomainObject) taskGroup.DomainObject {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) Save(do taskGroup.DomainObject) taskGroup.DomainObject {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) Delete(taskGroupId int) {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) SyncToData() {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) GetCanSchedulerTaskGroup(jobsName []string, ts time.Duration, count int, client vo.ClientVO) []vo.TaskDO {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) ToUnRunCount() int {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) ToSchedulerWorkingList() []taskGroup.DomainObject {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) ToFinishPageList(pageSize int, pageIndex int) []vo.TaskDO {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) GetTaskUnFinishList(jobsName []string, top int) []taskGroup.DomainObject {
	//TODO implement me
	panic("implement me")
}

func (repository taskGroupRepository) GetEnableTaskList(status eumTaskType.Enum, pageSize int, pageIndex int) core.PageList[vo.TaskDO] {
	//TODO implement me
	panic("implement me")
}
