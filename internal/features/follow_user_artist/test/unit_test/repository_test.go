package test_follow_user_artist

import (
	"errors"
	"testing"

	database "followservice/internal/db"
	mock_database "followservice/internal/db/test/mock"
	"followservice/internal/features/follow_user_artist"
	model "followservice/internal/model/domain"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var dbClient *mock_database.MockDatabaseClient
var followUserArtistRepository *follow_user_artist.FollowUserArtistRepository

func setUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbClient = mock_database.NewMockDatabaseClient(ctrl)
	followUserArtistRepository = follow_user_artist.NewFollowUserArtistRepository(database.NewDatabase(dbClient))
}

func TestAddNewPostMetaDataInRepository_WhenItReturnsSuccess(t *testing.T) {
	setUp(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	dbClient.EXPECT().CreateRelationship(newUserPair).Return(errors.New("some error"))

	err := followUserArtistRepository.AddUserRelationship(newUserPair)

	assert.NotNil(t, err)
}

func TestErrorOnAddNewPostMetaDataInRepository_WhenDatabaseFails(t *testing.T) {
	setUp(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	dbClient.EXPECT().CreateRelationship(newUserPair).Return(nil)

	err := followUserArtistRepository.AddUserRelationship(newUserPair)

	assert.Nil(t, err)
}
