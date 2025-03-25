package integration_test_assert

import (
	database "followservice/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertCachedUserFollowersExists(t *testing.T, db *database.Cache, username, lastFollowerId string, limit int, expectedFollowers []string) {
	cachedFollowers, cachedLastFollowerId, found := db.Client.GetUserFollowers(username, lastFollowerId, limit)
	assert.Equal(t, true, found)
	assert.Equal(t, expectedFollowers, cachedFollowers)
	assert.Equal(t, expectedFollowers[len(expectedFollowers)-1], cachedLastFollowerId)
}
