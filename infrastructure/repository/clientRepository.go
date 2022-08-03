package repository

import (
	"fss/domain/clients/client"
	"github.com/farseernet/farseer.go/cache/redis"
	"github.com/farseernet/farseer.go/core/container"
	"strconv"
)

func init() {
	// 注册仓储
	_ = container.Register(func() client.Repository {
		return &clientRepository{
			Client: redis.NewClient("default"),
		}
	})
}

const clientCacheKey = "FSS_ClientList"

type clientRepository struct {
	*redis.Client
}

func (repository clientRepository) ToList() []client.DomainObject {
	var clients []client.DomainObject
	_ = repository.Client.Hash.ToList(clientCacheKey, &clients)
	return clients
}

func (repository clientRepository) RemoveClient(id int64) {
	_, _ = repository.Client.Hash.Del(clientCacheKey, strconv.FormatInt(id, 10))
}

func (repository clientRepository) ToEntity(clientId int64) client.DomainObject {
	var client client.DomainObject
	_ = repository.Client.Hash.ToEntity(clientCacheKey, strconv.FormatInt(clientId, 10), &client)
	return client
}

func (repository clientRepository) Update(do client.DomainObject) client.DomainObject {
	var client client.DomainObject
	_ = repository.Client.Hash.SetEntity(clientCacheKey, strconv.FormatInt(do.Id, 10), &client)
	return client
}

func (repository clientRepository) GetCount() int64 {
	count := repository.Client.Hash.Count(clientCacheKey)
	return int64(count)
}
