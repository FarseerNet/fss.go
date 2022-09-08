package main

import (
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fss"
	"time"
)

// RegisterDemo3 Job3
func RegisterDemo3() {
	fss.RegisterJob("demo.job3", job3)
}
func job3(context fss.IFssContext) bool {
	task := context.GetTask()
	s := time.Now().Sub(task.StartAt)
	milliseconds := s.Milliseconds()
	if milliseconds > 0 {
		flog.Warningf("延迟了 %d ms", milliseconds)
	}
	flog.Printf("id=%d,caption=%s\n", task.TaskGroupId, task.Caption)
	return true
}
