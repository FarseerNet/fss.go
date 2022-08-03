package client

import "time"

type DomainObject struct {
	// 客户端ID
	Id int64
	// 客户端IP
	Ip string
	// 客户端名称
	Name string
	// 客户端能执行的任务
	Jobs []string
	// 活动时间
	ActivateAt time.Time
}

// IsTimeout 是否超时下线
func (do DomainObject) IsTimeout() bool {
	diffS := time.Now().Unix() - do.ActivateAt.Unix()
	return diffS >= 60
}
