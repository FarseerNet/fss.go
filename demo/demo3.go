package main

import (
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fss"
)

// RegisterDemo3 Job3
func RegisterDemo3() {
	fss.RegisterJob("demo.job3", job3)
}
func job3(context fss.IFssContext) bool {
	task := context.GetTask()
	flog.Printf("id=%d,caption=%s\n", task.TaskGroupId, task.Caption)
	return true
}
