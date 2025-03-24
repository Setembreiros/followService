package unit_test_get_user_followers

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"followservice/internal/features/get_user_followers"
	mock_get_user_followers "followservice/internal/features/get_user_followers/test/mock"

	"github.com/go-playground/assert/v2"
)

var controllerService *mock_get_user_followers.MockService
var controller *get_user_followers.GetUserFollowersController

func setUpHandler(t *testing.T) {
	setUp(t)
	controllerService = mock_get_user_followers.NewMockService(ctrl)
	controller = get_user_followers.NewGetUserFollowersController(controllerService)
}

func TestGetUserFollowersWithController_WhenSuccess(t *testing.T) {
	setUpHandler(t)
	expectedUsername := "usernameA"
	expectedLastFollowerId := "follower4"
	expectedLimit := 4
	req, _ := http.NewRequest("GET", fmt.Sprintf("/followers?username=%s&lastFollowerId=%s&limit=%d", expectedUsername, expectedLastFollowerId, expectedLimit), nil)
	ginContext.Request = req
	controllerService.EXPECT().GetUserFollowers(expectedUsername, expectedLastFollowerId, expectedLimit).Return([]string{"follower5", "follower6", "follower7"}, "follower7", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"followers":["follower5","follower6","follower7"],
			"lastFollowerId":"follower7"
		}
	}`

	controller.GetUserFollowers(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestGetUserFollowersWithController_WhenSuccessWithDefaultPaginationParameters(t *testing.T) {
	setUpHandler(t)
	expectedUsername := "usernameA"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/followers?username=%s", expectedUsername), nil)
	ginContext.Request = req
	expectedDefaultLastFollowerId := ""
	expectedDefaultLimit := 12
	controllerService.EXPECT().GetUserFollowers("usernameA", expectedDefaultLastFollowerId, expectedDefaultLimit).Return([]string{"follower5", "follower6", "follower7"}, "follower7", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"followers":["follower5","follower6","follower7"],
			"lastFollowerId":"follower7"
		}
	}`

	controller.GetUserFollowers(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnGetUserFollowersWithController_WhenServiceCallFails(t *testing.T) {
	setUpHandler(t)
	expectedUsername := "usernameA"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/followers?username=%s", expectedUsername), nil)
	ginContext.Request = req
	expectedError := errors.New("some error")
	controllerService.EXPECT().GetUserFollowers("usernameA", "", 12).Return([]string{}, "", expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.GetUserFollowers(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequestErrorOnGetUserPostsWithController_WhenLimitSmallerThanOne(t *testing.T) {
	setUpHandler(t)
	username := "usernameA"
	lastFollowerId := "follower4"
	wrongLimit := 0
	req, _ := http.NewRequest("GET", fmt.Sprintf("/followers?username=%s&lastFollowerId=%s&limit=%d", username, lastFollowerId, wrongLimit), nil)
	ginContext.Request = req
	expectedError := "Invalid pagination parameters, limit has to be greater than 0"
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError + `",
		"content":null
	}`

	controller.GetUserFollowers(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}
