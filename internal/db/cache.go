package database

//go:generate mockgen -source=cache.go -destination=test/mock/cache.go

type Cache struct {
	Client CacheClient
}

type CacheClient interface {
}

func NewCache(client CacheClient) *Cache {
	return &Cache{
		Client: client,
	}
}
