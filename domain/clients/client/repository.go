package client

import "github.com/farseer-go/collections"

type Repository interface {
	// ToList 获取客户端列表
	ToList() collections.List[DomainObject]
	// RemoveClient 移除客户端
	RemoveClient(id int64)
	// ToEntity 获取客户端
	ToEntity(clientId int64) DomainObject
	// Update 更新客户端的使用时间
	Update(do DomainObject)
	// GetCount 客户端数量
	GetCount() int64
}
