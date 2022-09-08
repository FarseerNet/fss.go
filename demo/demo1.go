package main

import (
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fss"
	"time"
)

// RegisterDemo1 Job1
func RegisterDemo1() {
	fss.RegisterJob("demo.job1", job1)
}
func job1(context fss.IFssContext) bool {
	task := context.GetTask()
	s := time.Now().Sub(task.StartAt)
	milliseconds := s.Milliseconds()
	if milliseconds > 0 {
		flog.Warningf("延迟了 %d ms", milliseconds)
	}
	flog.Printf("id=%d,caption=%s\n", task.TaskGroupId, task.Caption)
	return true
}
