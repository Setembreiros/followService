package database

//go:generate mockgen -source=cache.go -destination=test/mock/cache.go

type Cache struct {
	Client CacheClient
}

type CacheClient interface {
	Clean()
	SetUserFollowers(username string, lastFollowerId string, limit int, followers []string)
	GetUserFollowers(username string, lastFollowerId string, limit int) ([]string, string, bool)
}

func NewCache(client CacheClient) *Cache {
	return &Cache{
		Client: client,
	}
}
