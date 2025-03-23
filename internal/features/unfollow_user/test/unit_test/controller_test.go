package unit_test_ununfollow_user

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"followservice/internal/bus"
	"followservice/internal/features/unfollow_user"
	mock_unfollow_user "followservice/internal/features/unfollow_user/test/mock"
	model "followservice/internal/model/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

var controllerLoggerOutput bytes.Buffer
var controllerService *mock_unfollow_user.MockService
var controllerBus *bus.EventBus
var controller *unfollow_user.UnfollowUserController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	controllerService = mock_unfollow_user.NewMockService(ctrl)
	controllerBus = &bus.EventBus{}
	log.Logger = log.Output(&controllerLoggerOutput)
	controller = unfollow_user.NewUnfollowUserController(controllerService)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func TestUnfollowUser_WhenSuccess(t *testing.T) {
	setUpHandler(t)
	req, _ := http.NewRequest("DELETE", "/follow?followerId=usernameA&followeeId=usernameB", nil)
	ginContext.Request = req
	expectedUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	controllerService.EXPECT().UnfollowUser(expectedUserPair).Return(nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`

	controller.UnfollowUser(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequesErrortUnfollowUser_WhenMissingFollowerID(t *testing.T) {
	setUpHandler(t)
	req, _ := http.NewRequest("DELETE", "/follow?followeeId=usernameB", nil)
	ginContext.Request = req
	expectedBodyResponse := `{
		"error": true,
		"message": "Missing followerId parameter",
		"content": null
	}`

	controller.UnfollowUser(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestBadRequesErrortUnfollowUser_WhenMissingFolloweeID(t *testing.T) {
	setUpHandler(t)
	req, _ := http.NewRequest("DELETE", "/follow?followerId=usernameA", nil)
	ginContext.Request = req
	expectedBodyResponse := `{
		"error": true,
		"message": "Missing followeeId parameter",
		"content": null
	}`

	controller.UnfollowUser(ginContext)

	assert.Equal(t, apiResponse.Code, 400)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnUnfollowUser(t *testing.T) {
	setUpHandler(t)
	req, _ := http.NewRequest("DELETE", "/follow?followerId=usernameA&followeeId=usernameB", nil)
	ginContext.Request = req
	expectedUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	expectedError := errors.New("some error")
	controllerService.EXPECT().UnfollowUser(expectedUserPair).Return(expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.UnfollowUser(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
