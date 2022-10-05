package client

import "time"

type DomainObject struct {
	Id         int64     // 客户端ID
	Ip         string    // 客户端IP
	Name       string    // 客户端名称
	Jobs       []string  // 客户端能执行的任务
	ActivateAt time.Time // 活动时间
}

// IsTimeout 是否超时下线
func (do DomainObject) IsTimeout() bool {
	diffS := time.Now().Unix() - do.ActivateAt.Unix()
	return diffS >= 60
}
