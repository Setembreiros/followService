package integration_test_arrange

import (
	"context"
	"followservice/cmd/provider"
	database "followservice/internal/db"
	"testing"
)

func CreateTestCache(t *testing.T, ctx context.Context) *database.Cache {
	provider := provider.NewProvider("test")
	return provider.ProvideCache(ctx)
}

func AddCachedFollowersToCache(t *testing.T, cache *database.Cache, username, lastFollowerId string, limit int, followers []string) {
	cache.Client.SetUserFollowers(username, lastFollowerId, limit, followers)
}
