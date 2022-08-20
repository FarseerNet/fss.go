package repository

import (
	"fss/domain/clients/client"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/redis"
	"strconv"
)

func RegisterClientRepository() {
	// 注册仓储
	container.Register(func() client.Repository {
		return &clientRepository{
			Client: redis.NewClient("default"),
		}
	})
}

const clientCacheKey = "FSS_ClientList"

type clientRepository struct {
	*redis.Client
}

func (repository clientRepository) ToList() collections.List[client.DomainObject] {
	var clients []client.DomainObject
	_ = repository.Client.Hash.ToArray(clientCacheKey, &clients)
	return collections.NewList(clients...)
}

func (repository clientRepository) RemoveClient(id int64) {
	_, _ = repository.Client.Hash.Del(clientCacheKey, strconv.FormatInt(id, 10))
}

func (repository clientRepository) ToEntity(clientId int64) client.DomainObject {
	var do client.DomainObject
	_ = repository.Client.Hash.ToEntity(clientCacheKey, strconv.FormatInt(clientId, 10), &do)
	return do
}

func (repository clientRepository) Update(do client.DomainObject) {
	_ = repository.Client.Hash.SetEntity(clientCacheKey, strconv.FormatInt(do.Id, 10), &do)
}

func (repository clientRepository) GetCount() int64 {
	count := repository.Client.Hash.Count(clientCacheKey)
	return int64(count)
}
