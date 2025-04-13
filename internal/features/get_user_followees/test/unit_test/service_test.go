package unit_test_get_user_followees

import (
	"errors"
	"testing"

	"followservice/internal/features/get_user_followees"
	mock_get_user_followees "followservice/internal/features/get_user_followees/test/mock"

	"github.com/stretchr/testify/assert"
)

var serviceRepository *mock_get_user_followees.MockRepository
var service *get_user_followees.GetUserFolloweesService

func setUpService(t *testing.T) {
	setUp(t)
	serviceRepository = mock_get_user_followees.NewMockRepository(ctrl)
	service = get_user_followees.NewGetUserFolloweesService(serviceRepository)
}

func TestGetUserFolloweesWithService_WhenSuccess(t *testing.T) {
	setUpService(t)
	username := "usernameA"
	lastFolloweeId := "followee4"
	limit := 4
	expectedFollowees := []string{"followee5", "followee6", "followee7"}
	expectedLastFolloweeId := "followee7"
	serviceRepository.EXPECT().GetUserFollowees(username, lastFolloweeId, limit).Return(expectedFollowees, expectedLastFolloweeId, nil)

	followees, lastFolloweeId, err := service.GetUserFollowees(username, lastFolloweeId, limit)

	assert.Nil(t, err)
	assert.Equal(t, followees, expectedFollowees)
	assert.Equal(t, lastFolloweeId, expectedLastFolloweeId)
}

func TestErrorOnGetUserFolloweesWithService_WhenGetUserFolloweesFails(t *testing.T) {
	setUpService(t)
	username := "usernameA"
	lastFolloweeId := "followee4"
	limit := 4
	expectedFollowees := []string{}
	expectedLastFolloweeId := ""
	serviceRepository.EXPECT().GetUserFollowees(username, lastFolloweeId, limit).Return(expectedFollowees, expectedLastFolloweeId, errors.New("some error"))

	followees, lastFolloweeId, err := service.GetUserFollowees(username, lastFolloweeId, limit)

	assert.NotNil(t, err)
	assert.Equal(t, followees, expectedFollowees)
	assert.Equal(t, lastFolloweeId, expectedLastFolloweeId)
	assert.Contains(t, loggerOutput.String(), "Error getting  "+username+"'s followees")
}
