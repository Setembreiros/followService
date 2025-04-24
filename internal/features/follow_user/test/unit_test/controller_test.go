package unit_test_follow_user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"followservice/internal/bus"
	"followservice/internal/features/follow_user"
	mock_follow_user "followservice/internal/features/follow_user/test/mock"
	model "followservice/internal/model/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

var controllerLoggerOutput bytes.Buffer
var controllerService *mock_follow_user.MockService
var controllerBus *bus.EventBus
var controller *follow_user.FollowUserController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	controllerService = mock_follow_user.NewMockService(ctrl)
	controllerBus = &bus.EventBus{}
	log.Logger = log.Output(&controllerLoggerOutput)
	controller = follow_user.NewFollowUserController(controllerService)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func TestFollowUser(t *testing.T) {
	setUpHandler(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	data, _ := serializeData(newUserPair)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/follow", bytes.NewBuffer(data))
	controllerService.EXPECT().FollowUser(newUserPair).Return(nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`

	controller.FollowUser(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerErrorOnFollowUser(t *testing.T) {
	setUpHandler(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	data, _ := serializeData(newUserPair)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/follow", bytes.NewBuffer(data))
	expectedError := errors.New("some error")
	controllerService.EXPECT().FollowUser(newUserPair).Return(expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content": null
	}`

	controller.FollowUser(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func serializeData(data any) ([]byte, error) {
	return json.Marshal(data)
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
