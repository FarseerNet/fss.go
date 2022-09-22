package interfaces

import (
	"fss/application/clients/clientApp"
	"fss/application/log/taskLogApp"
	logReq "fss/application/log/taskLogApp/request"
	"fss/application/tasks/taskGroupApp"
	taskReq "fss/application/tasks/taskGroupApp/request"
	"github.com/beego/beego/v2/server/web"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/stopwatch"
)

type MetaController struct {
	web.Controller
}

// GetClientList 取出全局客户端列表
func (r *MetaController) GetClientList() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "GetClientList，耗时：%s", sw.GetMillisecondsText())
	}()

	//调用应用层
	result := clientApp.ToList()
	apiResponse := core.Success("请求成功", result)
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// GetClientCount 客户端数量
func (r *MetaController) GetClientCount() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "GetClientCount，耗时：%s", sw.GetMillisecondsText())
	}()

	//调用应用层
	result := clientApp.GetCount()
	apiResponse := core.ApiResponseLongSuccess("请求成功", result)

	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// GetRunLogList 获取日志
func (r *MetaController) GetRunLogList() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "GetRunLogList，耗时：%s", sw.GetMillisecondsText())
	}()

	// 读取结构数据
	var req logReq.GetRunLogRequest
	_ = r.BindJSON(&req)
	//调用应用层
	result := taskLogApp.GetList(req)
	apiResponse := core.Success("请求成功", result.ToArray())
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// CopyTaskGroup 复制任务组
func (r *MetaController) CopyTaskGroup() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "CopyTaskGroup，耗时：%s", sw.GetMillisecondsText())
	}()

	//读取结构数据
	var req taskReq.OnlyIdRequest
	_ = r.BindJSON(&req)
	//调用应用层
	result := taskGroupApp.CopyTaskGroup(req)
	apiResponse := core.ApiResponseIntSuccess("请求成功", result)
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// DeleteTaskGroup 删除任务组
func (r *MetaController) DeleteTaskGroup() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "DeleteTaskGroup，耗时：%s", sw.GetMillisecondsText())
	}()

	//读取结构数据
	var req taskReq.OnlyIdRequest
	_ = r.BindJSON(&req)
	//调用应用层
	taskGroupApp.Delete(req)
	apiResponse := core.ApiResponseStringSuccess("请求成功")
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// AddTaskGroup 添加任务组信息
func (r *MetaController) AddTaskGroup() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "AddTaskGroup，耗时：%s", sw.GetMillisecondsText())
	}()

	//读取结构数据
	var req taskGroupApp.DTO
	_ = r.BindJSON(&req)
	//调用应用层
	result := taskGroupApp.Add(req)
	apiResponse := core.ApiResponseIntSuccess("请求成功", result)
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// SaveTaskGroup 保存任务组
func (r *MetaController) SaveTaskGroup() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "SaveTaskGroup，耗时：%s", sw.GetMillisecondsText())
	}()

	//读取结构数据
	var req taskGroupApp.DTO
	_ = r.BindJSON(&req)
	//调用应用层
	taskGroupApp.Save(req)
	apiResponse := core.ApiResponseStringSuccess("请求成功")
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// CancelTask 取消任务
func (r *MetaController) CancelTask() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "CancelTask，耗时：%s", sw.GetMillisecondsText())
	}()

	//读取结构数据
	var req taskReq.OnlyIdRequest
	_ = r.BindJSON(&req)
	//调用应用层
	taskGroupApp.CancelTask(req)
	apiResponse := core.ApiResponseStringSuccess("请求成功")
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// SyncCacheToDb 同步数据
func (r *MetaController) SyncCacheToDb() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "SyncCacheToDb，耗时：%s", sw.GetMillisecondsText())
	}()

	//调用应用层
	taskGroupApp.SyncTaskGroup()
	apiResponse := core.ApiResponseStringSuccess("请求成功")
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// GetTaskGroupInfo 获取任务组信息
func (r *MetaController) GetTaskGroupInfo() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "GetTaskGroupInfo，耗时：%s", sw.GetMillisecondsText())
	}()

	//读取结构数据
	var req taskReq.OnlyIdRequest
	_ = r.BindJSON(&req)
	//调用应用层
	result := taskGroupApp.ToEntity(req)
	if result.Id > 1 {
		apiResponse := core.Success("请求成功", result)
		// 响应数据
		r.Data["json"] = &apiResponse
	} else {
		apiResponse := core.Error403("数据不存在")
		// 响应数据
		r.Data["json"] = &apiResponse
	}

	_ = r.ServeJSON()
}

// GetTaskGroupList 获取所有任务组中的任务
func (r *MetaController) GetTaskGroupList() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "GetTaskGroupList，耗时：%s", sw.GetMillisecondsText())
	}()

	//调用应用层
	result := taskGroupApp.ToList()
	apiResponse := core.Success("请求成功", result)
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// TodayTaskFailCount 今日执行失败数量
func (r *MetaController) TodayTaskFailCount() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "TodayTaskFailCount，耗时：%s", sw.GetMillisecondsText())
	}()

	//调用应用层
	result := taskGroupApp.TodayFailCount()
	apiResponse := core.ApiResponseLongSuccess("请求成功", result)
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// GetTaskGroupCount 获取任务组数量
func (r *MetaController) GetTaskGroupCount() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "GetTaskGroupCount，耗时：%s", sw.GetMillisecondsText())
	}()

	//调用应用层
	result := taskGroupApp.GetTaskGroupCount()
	apiResponse := core.ApiResponseLongSuccess("请求成功", result)
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// GetTaskGroupUnRunCount 获取未执行的任务数量
func (r *MetaController) GetTaskGroupUnRunCount() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "GetTaskGroupUnRunCount，耗时：%s", sw.GetMillisecondsText())
	}()

	//调用应用层
	result := taskGroupApp.ToUnRunCount()
	apiResponse := core.ApiResponseIntSuccess("请求成功", result)
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// GetTaskList 获取指定任务组的任务列表（FOPS）
func (r *MetaController) GetTaskList() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "GetTaskList，耗时：%s", sw.GetMillisecondsText())
	}()

	//读取结构数据
	var req taskReq.GetTaskListRequest
	_ = r.BindJSON(&req)
	//调用应用层
	result := taskGroupApp.GetTaskList(req)
	apiResponse := core.Success("请求成功", result)
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// GetTaskFinishList 获取已完成的任务列表
func (r *MetaController) GetTaskFinishList() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "GetTaskFinishList，耗时：%s", sw.GetMillisecondsText())
	}()

	//读取结构数据
	var req taskReq.PageSizeRequest
	_ = r.BindJSON(&req)
	//调用应用层
	result := taskGroupApp.ToFinishPageList(req)
	apiResponse := core.Success("请求成功", result)
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// GetTaskUnFinishList 获取进行中的任务
func (r *MetaController) GetTaskUnFinishList() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "GetTaskUnFinishList，耗时：%s", sw.GetMillisecondsText())
	}()

	//读取结构数据
	var req taskReq.OnlyTopRequest
	_ = r.BindJSON(&req)
	//调用应用层
	result := taskGroupApp.GetTaskUnFinishList(req)
	apiResponse := core.Success("请求成功", result)
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// GetEnableTaskList 获取在用的任务组
func (r *MetaController) GetEnableTaskList() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "GetEnableTaskList，耗时：%s", sw.GetMillisecondsText())
	}()

	//读取结构数据
	var req taskReq.GetAllTaskListRequest
	_ = r.BindJSON(&req)
	//调用应用层
	result := taskGroupApp.GetEnableTaskList(req)
	apiResponse := core.Success("请求成功", result)
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}

// SetEnable 设置任务组的状态
func (r *MetaController) SetEnable() {
	sw := stopwatch.StartNew()
	defer func() {
		flog.ComponentInfof("webapi", "SetEnable，耗时：%s", sw.GetMillisecondsText())
	}()

	//读取结构数据
	var req taskReq.SetEnableRequest
	_ = r.BindJSON(&req)
	//调用应用层
	taskGroupApp.SetEnable(req)
	apiResponse := core.ApiResponseStringSuccess("请求成功")
	// 响应数据
	r.Data["json"] = &apiResponse
	_ = r.ServeJSON()
}
