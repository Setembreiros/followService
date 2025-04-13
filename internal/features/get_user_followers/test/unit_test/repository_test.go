package unit_test_get_user_followers

import (
	"errors"
	"testing"

	database "followservice/internal/db"
	mock_database "followservice/internal/db/test/mock"
	"followservice/internal/features/get_user_followers"

	"github.com/stretchr/testify/assert"
)

var dbClient *mock_database.MockDatabaseClient
var repository *get_user_followers.GetUserFollowersRepository

func setUpRepository(t *testing.T) {
	setUp(t)
	dbClient = mock_database.NewMockDatabaseClient(ctrl)
	repository = get_user_followers.NewGetUserFollowersRepository(database.NewDatabase(dbClient))
}

func TestGetUserFollowersFromRepository_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	lastFollowerId := "follower4"
	limit := 4
	expectedFollowers := []string{"follower5", "follower6", "follower7"}
	expectedLastFollowerId := "follower4"
	dbClient.EXPECT().GetUserFollowers(username, lastFollowerId, limit).Return(expectedFollowers, expectedLastFollowerId, nil)

	followers, lastFollowerId, err := repository.GetUserFollowers(username, lastFollowerId, limit)

	assert.Nil(t, err)
	assert.Equal(t, followers, expectedFollowers)
	assert.Equal(t, lastFollowerId, expectedLastFollowerId)
}

func TestErrorOnGetUserFollowersFromRepository_WhenDatabaseFails(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	lastFollowerId := "follower4"
	limit := 4
	dbClient.EXPECT().GetUserFollowers(username, lastFollowerId, limit).Return([]string{}, "", errors.New("some error"))

	followers, lastFollowerId, err := repository.GetUserFollowers(username, lastFollowerId, limit)

	assert.NotNil(t, err)
	assert.Equal(t, followers, []string{})
	assert.Equal(t, lastFollowerId, "")
}
