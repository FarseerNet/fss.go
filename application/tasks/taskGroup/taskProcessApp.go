package taskGroup

import (
	"fmt"
	"fss/application/clients/client"
	"fss/application/tasks/taskGroup/request"
	"fss/domain/_/eumTaskType"
	"fss/domain/log"
	"fss/domain/tasks/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/exception"
)

type taskProcessApp struct {
	repository taskGroup.Repository
}

func NewTaskProcessApp() *taskProcessApp {
	return &taskProcessApp{repository: container.Resolve[taskGroup.Repository]()}
}

// JobInvoke 客户端执行任务
func (r *taskProcessApp) JobInvoke(dto request.JobInvokeDTO) string {
	clientDTO := client.NewApp().GetClient()
	taskGroupDO := r.repository.ToEntity(dto.TaskGroupId)

	if taskGroupDO.Id < 1 {
		taskGroupNotExistsMsg := fmt.Sprintf("所属的任务组：%d 不存在", dto.TaskGroupId)
		log.TaskLogAddService(dto.TaskGroupId, "", "", eumLogLevel.Warning, taskGroupNotExistsMsg)
		exception.ThrowRefuseException(taskGroupNotExistsMsg)
	}

	defer exception.Catch().
		RefuseException(func(exp *exception.RefuseException) {
			if taskGroupDO.Id > 0 {
				log.TaskLogAddService(dto.TaskGroupId, taskGroupDO.JobName, taskGroupDO.Caption, eumLogLevel.Warning, exp.Message)
			}
			exp.ContinueRecover(exp.Message)
		}).
		String(func(exp string) {
			if taskGroupDO.Id > 0 {
				taskGroupDO.Cancel()
				r.repository.Save(taskGroupDO)
				log.TaskLogAddService(taskGroupDO.Id, taskGroupDO.JobName, taskGroupDO.Caption, eumLogLevel.Error, exp)
			}
		})

	// 如果有日志
	if dto.Log.Log != "" {
		log.TaskLogAddService(taskGroupDO.Id, taskGroupDO.JobName, taskGroupDO.Caption, dto.Log.LogLevel, dto.Log.Log)
	}

	// 不相等，说明被覆盖了（JOB请求慢了。被调度重新执行了）
	if taskGroupDO.Task.Client.Id > 0 && taskGroupDO.Task.Client.Id != clientDTO.Id {
		exception.ThrowRefuseExceptionf("任务： %s（%s） ，{%d}与本次请求的客户端{%d} 不一致，忽略本次请求", taskGroupDO.Caption, taskGroupDO.JobName, taskGroupDO.Task.Client.Id, clientDTO.Id)
	}

	// 更新执行中状态
	taskGroupDO.Working(dto.Data, dto.NextTimespan, dto.Progress, dto.Status, dto.RunSpeed)
	r.repository.Save(taskGroupDO)

	if dto.Status != eumTaskType.Working && dto.Status != eumTaskType.Success {
		exception.ThrowRefuseExceptionf("任务组：TaskGroupId=%d，Caption=%s，JobName=%s 执行失败", taskGroupDO.Id, taskGroupDO.Caption, taskGroupDO.JobName)
	}
	return fmt.Sprintf("任务组：TaskGroupId=%d，Caption=%s，JobName=%s 处理成功", dto.TaskGroupId, taskGroupDO.Caption, taskGroupDO.JobName)

}
