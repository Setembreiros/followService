package integration_test_arrange

import (
	"context"
	"followservice/cmd/provider"
	database "followservice/internal/db"
	model "followservice/internal/model/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateTestDatabase(t *testing.T, ctx context.Context) *database.Database {
	provider := provider.NewProvider("test")
	return provider.ProvideDb(ctx)
}

func AddRelationshipToDatabase(t *testing.T, db *database.Database, userPair *model.UserPairRelationship) {
	err := db.Client.CreateRelationship(userPair)
	assert.Nil(t, err)
	existsInDatabase, err := db.Client.RelationshipExists(userPair)
	assert.Nil(t, err)
	assert.Equal(t, existsInDatabase, true)
}

func PopulateDb(t *testing.T, db *database.Database, followeeId, lastFollowerId string) {
	existingUserPairs := []*model.UserPairRelationship{
		{
			FollowerID: lastFollowerId,
			FolloweeID: followeeId,
		},
		{
			FollowerID: "username1",
			FolloweeID: followeeId,
		},
		{
			FollowerID: "username2",
			FolloweeID: followeeId,
		},
		{
			FollowerID: "username3",
			FolloweeID: followeeId,
		},
		{
			FollowerID: "username4",
			FolloweeID: followeeId,
		},
		{
			FollowerID: "username5",
			FolloweeID: followeeId,
		},
	}

	for _, existingUserPair := range existingUserPairs {
		AddRelationshipToDatabase(t, db, existingUserPair)
	}
}
