package unit_test_follow_user

import (
	"errors"
	"testing"

	database "followservice/internal/db"
	mock_database "followservice/internal/db/test/mock"
	"followservice/internal/features/follow_user"
	model "followservice/internal/model/domain"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var dbClient *mock_database.MockDatabaseClient
var followUserRepository *follow_user.FollowUserRepository

func setUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbClient = mock_database.NewMockDatabaseClient(ctrl)
	followUserRepository = follow_user.NewFollowUserRepository(database.NewDatabase(dbClient))
}

func TestAddNewPostMetaDataInRepository_WhenItReturnsSuccess(t *testing.T) {
	setUp(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	dbClient.EXPECT().CreateRelationship(newUserPair).Return(errors.New("some error"))

	err := followUserRepository.AddUserRelationship(newUserPair)

	assert.NotNil(t, err)
}

func TestErrorOnAddNewPostMetaDataInRepository_WhenDatabaseFails(t *testing.T) {
	setUp(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	dbClient.EXPECT().CreateRelationship(newUserPair).Return(nil)

	err := followUserRepository.AddUserRelationship(newUserPair)

	assert.Nil(t, err)
}
