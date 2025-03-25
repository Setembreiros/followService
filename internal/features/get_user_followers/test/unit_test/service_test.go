package unit_test_get_user_followers

import (
	"errors"
	"testing"

	"followservice/internal/features/get_user_followers"
	mock_get_user_followers "followservice/internal/features/get_user_followers/test/mock"

	"github.com/stretchr/testify/assert"
)

var serviceRepository *mock_get_user_followers.MockRepository
var service *get_user_followers.GetUserFollowersService

func setUpService(t *testing.T) {
	setUp(t)
	serviceRepository = mock_get_user_followers.NewMockRepository(ctrl)
	service = get_user_followers.NewGetUserFollowersService(serviceRepository)
}

func TestGetUserFollowersWithService_WhenSuccess(t *testing.T) {
	setUpService(t)
	username := "usernameA"
	lastFollowerId := "follower4"
	limit := 4
	expectedFollowers := []string{"follower5", "follower6", "follower7"}
	expectedLastFollowerId := "follower7"
	serviceRepository.EXPECT().GetUserFollowers(username, lastFollowerId, limit).Return(expectedFollowers, expectedLastFollowerId, nil)

	followers, lastFollowerId, err := service.GetUserFollowers(username, lastFollowerId, limit)

	assert.Nil(t, err)
	assert.Equal(t, followers, expectedFollowers)
	assert.Equal(t, lastFollowerId, expectedLastFollowerId)
}

func TestErrorOnGetUserFollowersWithService_WhenGetUserFollowersFails(t *testing.T) {
	setUpService(t)
	username := "usernameA"
	lastFollowerId := "follower4"
	limit := 4
	expectedFollowers := []string{}
	expectedLastFollowerId := ""
	serviceRepository.EXPECT().GetUserFollowers(username, lastFollowerId, limit).Return(expectedFollowers, expectedLastFollowerId, errors.New("some error"))

	followers, lastFollowerId, err := service.GetUserFollowers(username, lastFollowerId, limit)

	assert.NotNil(t, err)
	assert.Equal(t, followers, expectedFollowers)
	assert.Equal(t, lastFollowerId, expectedLastFollowerId)
	assert.Contains(t, loggerOutput.String(), "Error getting  "+username+"'s followers")
}
