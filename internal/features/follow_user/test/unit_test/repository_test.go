package unit_test_follow_user

import (
	"bytes"
	"errors"
	"testing"

	database "followservice/internal/db"
	mock_database "followservice/internal/db/test/mock"
	"followservice/internal/features/follow_user"
	model "followservice/internal/model/domain"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var dbClient *mock_database.MockDatabaseClient
var followUserRepository *follow_user.FollowUserRepository
var repositoryLoggerOutput bytes.Buffer

func setUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	log.Logger = log.Output(&repositoryLoggerOutput)
	dbClient = mock_database.NewMockDatabaseClient(ctrl)
	followUserRepository = follow_user.NewFollowUserRepository(database.NewDatabase(dbClient))
}

func TestAddUserRelationshipInRepository_WhenItReturnsSuccess(t *testing.T) {
	setUp(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	dbClient.EXPECT().RelationshipExists(newUserPair).Return(false, nil)
	dbClient.EXPECT().CreateRelationship(newUserPair).Return(nil)

	err := followUserRepository.AddUserRelationship(newUserPair)

	assert.Nil(t, err)
	assert.NotContains(t, repositoryLoggerOutput.String(), "Relationship already exists, "+newUserPair.FollowerID+" -> "+newUserPair.FolloweeID)
}

func TestErrorOnAddUserRelationshipInRepository_WhenCreateRelationshipFails(t *testing.T) {
	setUp(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	dbClient.EXPECT().RelationshipExists(newUserPair).Return(false, nil)
	dbClient.EXPECT().CreateRelationship(newUserPair).Return(errors.New("some error"))

	err := followUserRepository.AddUserRelationship(newUserPair)

	assert.NotNil(t, err)
}

func TestErrorOnAddUserRelationshipInRepository_WhenRelationshipExistsFails(t *testing.T) {
	setUp(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	dbClient.EXPECT().RelationshipExists(newUserPair).Return(false, errors.New("some error"))

	err := followUserRepository.AddUserRelationship(newUserPair)

	assert.NotNil(t, err)
	assert.NotContains(t, repositoryLoggerOutput.String(), "Relationship already exists, "+newUserPair.FollowerID+" -> "+newUserPair.FolloweeID)
}

func TesErrorOnAddUserRelationshipInRepository_WhenRelationshipAlreadyExists(t *testing.T) {
	setUp(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	dbClient.EXPECT().RelationshipExists(newUserPair).Return(true, nil)

	err := followUserRepository.AddUserRelationship(newUserPair)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Relationship already exists, "+newUserPair.FollowerID+" -> "+newUserPair.FolloweeID)
}
