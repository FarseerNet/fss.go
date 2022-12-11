package main

import (
	"fss/application/clients/clientApp"
	"fss/application/log/taskLogApp"
	"fss/application/tasks/taskGroupApp"
	"github.com/farseer-go/webapi"
)

var routeMeta = []webapi.Route{
	{Url: "/meta/GetClientList", Method: "POST", Action: clientApp.ToList},
	{Url: "/meta/GetClientCount", Method: "POST", Action: clientApp.GetCount},
	{Url: "/meta/GetRunLogList", Method: "POST", Action: taskLogApp.GetList},
	{Url: "/meta/CopyTaskGroup", Method: "POST", Action: taskGroupApp.CopyTaskGroup},
	{Url: "/meta/DeleteTaskGroup", Method: "POST", Action: taskGroupApp.Delete},
	{Url: "/meta/AddTaskGroup", Method: "POST", Action: taskGroupApp.Add},
	{Url: "/meta/SaveTaskGroup", Method: "POST", Action: taskGroupApp.Save},
	{Url: "/meta/CancelTask", Method: "POST", Action: taskGroupApp.CancelTask},
	{Url: "/meta/SyncCacheToDb", Method: "POST", Action: taskGroupApp.SyncTaskGroup},
	{Url: "/meta/GetTaskGroupInfo", Method: "POST", Action: taskGroupApp.ToEntity},
	{Url: "/meta/TodayTaskFailCount", Method: "POST", Action: taskGroupApp.TodayTaskFailCount},
	{Url: "/meta/GetTaskGroupCount", Method: "POST", Action: taskGroupApp.GetTaskGroupCount},
	{Url: "/meta/GetTaskGroupUnRunCount", Method: "POST", Action: taskGroupApp.GetTaskGroupUnRunCount},
	{Url: "/meta/GetTaskList", Method: "POST", Action: taskGroupApp.GetTaskList},
	{Url: "/meta/GetTaskFinishList", Method: "POST", Action: taskGroupApp.GetTaskFinishList},
	{Url: "/meta/GetTaskUnFinishList", Method: "POST", Action: taskGroupApp.GetTaskUnFinishList},
	{Url: "/meta/GetEnableTaskList", Method: "POST", Action: taskGroupApp.GetEnableTaskList},
	{Url: "/meta/SetEnable", Method: "POST", Action: taskGroupApp.SetEnable},
	{Url: "/meta/GetTaskGroupList", Method: "POST", Action: taskGroupApp.ToList},
}
