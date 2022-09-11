package main

import (
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fss"
)

// RegisterDemo2 Job2
func RegisterDemo2() {
	fss.RegisterJob("demo.job2", job2)
}
func job2(context fss.IFssContext) bool {
	task := context.GetTask()
	flog.Printf("id=%d,caption=%s\n", task.TaskGroupId, task.Caption)
	return true
}
