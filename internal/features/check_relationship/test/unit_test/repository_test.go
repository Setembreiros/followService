package unit_test_check_relationship

import (
	"errors"
	"testing"

	database "followservice/internal/db"
	mock_database "followservice/internal/db/test/mock"
	"followservice/internal/features/check_relationship"
	model "followservice/internal/model/domain"

	"github.com/stretchr/testify/assert"
)

var dbClient *mock_database.MockDatabaseClient
var repository *check_relationship.CheckRelationshipRepository

func setUpRepository(t *testing.T) {
	setUp(t)
	dbClient = mock_database.NewMockDatabaseClient(ctrl)
	repository = check_relationship.NewCheckRelationshipRepository(database.NewDatabase(dbClient))
}

func TestCheckRelationshipRepository_WhenExists(t *testing.T) {
	setUpRepository(t)
	userPair := &model.UserPairRelationship{FollowerID: "userA", FolloweeID: "userB"}
	dbClient.EXPECT().RelationshipExists(userPair).Return(true, nil)

	result, err := repository.CheckRelationship(userPair)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestCheckRelationshipRepository_WhenNotExists(t *testing.T) {
	setUpRepository(t)
	userPair := &model.UserPairRelationship{FollowerID: "userA", FolloweeID: "userB"}
	dbClient.EXPECT().RelationshipExists(userPair).Return(false, nil)

	result, err := repository.CheckRelationship(userPair)

	assert.Nil(t, err)
	assert.False(t, result)
}

func TestCheckRelationshipRepository_WhenError(t *testing.T) {
	setUpRepository(t)
	userPair := &model.UserPairRelationship{FollowerID: "userA", FolloweeID: "userB"}
	dbClient.EXPECT().RelationshipExists(userPair).Return(false, errors.New("db error"))

	result, err := repository.CheckRelationship(userPair)

	assert.NotNil(t, err)
	assert.False(t, result)
}
