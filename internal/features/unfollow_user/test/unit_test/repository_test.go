package unit_test_unfollow_user

import (
	"errors"
	"testing"

	database "followservice/internal/db"
	mock_database "followservice/internal/db/test/mock"
	"followservice/internal/features/unfollow_user"
	model "followservice/internal/model/domain"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var dbClient *mock_database.MockDatabaseClient
var unfollowUserRepository *unfollow_user.UnfollowUserRepository

func setUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbClient = mock_database.NewMockDatabaseClient(ctrl)
	unfollowUserRepository = unfollow_user.NewUnfollowUserRepository(database.NewDatabase(dbClient))
}

func TestRemoveUserRelationshipInRepository_WhenItReturnsSuccess(t *testing.T) {
	setUp(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	dbClient.EXPECT().DeleteRelationship(newUserPair).Return(nil)

	err := unfollowUserRepository.RemoveUserRelationship(newUserPair)

	assert.Nil(t, err)
}

func TestErrorOnRemoveUserRelationshipInRepository_WhenDatabaseFails(t *testing.T) {
	setUp(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	dbClient.EXPECT().DeleteRelationship(newUserPair).Return(errors.New("some error"))

	err := unfollowUserRepository.RemoveUserRelationship(newUserPair)

	assert.NotNil(t, err)
}
