package taskGroupApp

import (
	"fmt"
	"fss/application/clients/clientApp"
	"fss/application/tasks/taskGroupApp/request"
	"fss/domain/_/eumTaskType"
	"fss/domain/log"
	"fss/domain/tasks/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/exception"
)

// JobInvoke 客户端执行任务
func JobInvoke(clientDTO clientApp.DTO, dto request.JobInvokeDTO) string {
	repository := container.Resolve[taskGroup.Repository]()
	taskGroupDO := repository.ToEntity(dto.TaskGroupId)

	try := exception.Try(func() {
		if taskGroupDO.Id < 1 {
			taskGroupNotExistsMsg := fmt.Sprintf("所属的任务组：%d 不存在", dto.TaskGroupId)
			log.TaskLogAddService(dto.TaskGroupId, "", "", eumLogLevel.Warning, taskGroupNotExistsMsg)
			exception.ThrowRefuseException(taskGroupNotExistsMsg)
		}

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
		repository.Save(taskGroupDO)

		if dto.Status != eumTaskType.Working && dto.Status != eumTaskType.Success {
			exception.ThrowRefuseExceptionf("任务组：TaskGroupId=%d，Caption=%s，JobName=%s 执行失败", taskGroupDO.Id, taskGroupDO.Caption, taskGroupDO.JobName)
		}
	})

	try.CatchRefuseException(func(exp *exception.RefuseException) {
		if taskGroupDO.Id > 0 {
			log.TaskLogAddService(dto.TaskGroupId, taskGroupDO.JobName, taskGroupDO.Caption, eumLogLevel.Warning, exp.Message)
		}
		// 让api层处理
		exception.ThrowRefuseException(exp.Message)
	})

	try.CatchStringException(func(exp string) {
		if taskGroupDO.Id > 0 {
			taskGroupDO.Cancel()
			repository.Save(taskGroupDO)
			log.TaskLogAddService(taskGroupDO.Id, taskGroupDO.JobName, taskGroupDO.Caption, eumLogLevel.Error, exp)
		}
		// 让api层处理
		exception.ThrowException(exp)
	})

	try.ThrowUnCatch()

	return fmt.Sprintf("任务组：TaskGroupId=%d，Caption=%s，JobName=%s 处理成功", dto.TaskGroupId, taskGroupDO.Caption, taskGroupDO.JobName)
}
