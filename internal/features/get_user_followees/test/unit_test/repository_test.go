package unit_test_get_user_followees

import (
	"errors"
	"testing"

	database "followservice/internal/db"
	mock_database "followservice/internal/db/test/mock"
	"followservice/internal/features/get_user_followees"

	"github.com/stretchr/testify/assert"
)

var dbClient *mock_database.MockDatabaseClient
var repository *get_user_followees.GetUserFolloweesRepository

func setUpRepository(t *testing.T) {
	setUp(t)
	dbClient = mock_database.NewMockDatabaseClient(ctrl)
	repository = get_user_followees.NewGetUserFolloweesRepository(database.NewDatabase(dbClient))
}

func TestGetUserFolloweesFromRepository_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	lastFolloweeId := "followee4"
	limit := 4
	expectedFollowees := []string{"followee5", "followee6", "followee7"}
	expectedLastFolloweeId := "followee4"
	dbClient.EXPECT().GetUserFollowees(username, lastFolloweeId, limit).Return(expectedFollowees, expectedLastFolloweeId, nil)

	followees, lastFolloweeId, err := repository.GetUserFollowees(username, lastFolloweeId, limit)

	assert.Nil(t, err)
	assert.Equal(t, followees, expectedFollowees)
	assert.Equal(t, lastFolloweeId, expectedLastFolloweeId)
}

func TestErrorOnGetUserFolloweesFromRepository_WhenDatabaseFails(t *testing.T) {
	setUpRepository(t)
	username := "usernameA"
	lastFolloweeId := "followee4"
	limit := 4
	dbClient.EXPECT().GetUserFollowees(username, lastFolloweeId, limit).Return([]string{}, "", errors.New("some error"))

	followees, lastFolloweeId, err := repository.GetUserFollowees(username, lastFolloweeId, limit)

	assert.NotNil(t, err)
	assert.Equal(t, followees, []string{})
	assert.Equal(t, lastFolloweeId, "")
}
