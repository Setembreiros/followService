package integration_test_follow_user

import (
	"bytes"
	"encoding/json"
	"followservice/internal/bus"
	mock_bus "followservice/internal/bus/test/mock"
	database "followservice/internal/db"
	"followservice/internal/features/follow_user"
	model "followservice/internal/model/domain"
	"followservice/internal/model/events"
	integration_test_assert "followservice/test/integration_test_common/assert"
	integration_test_builder "followservice/test/integration_test_common/builder"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

var db *database.Database
var controller *follow_user.FollowUserController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context
var serviceExternalBus *mock_bus.MockExternalBus

func setUp(t *testing.T) {
	// Mocks
	ctrl := gomock.NewController(t)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus := bus.NewEventBus(serviceExternalBus)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)

	// Real infrastructure and services
	db = integration_test_builder.NewDatabaseBuilder(t, ginContext).Build()
	repository := follow_user.NewFollowUserRepository(db)
	service := follow_user.NewFollowUserService(repository, serviceBus)
	controller = follow_user.NewFollowUserController(service)
}

func tearDown() {
	db.Client.Clean()
}

func TestFollowUser_WhenItReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	data, _ := serializeData(newUserPair)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/follow", bytes.NewBuffer(data))
	expectedUserAFollowedUserBEvent := &events.UserAFollowedUserBEvent{
		FollowerID: newUserPair.FollowerID,
		FolloweeID: newUserPair.FolloweeID,
	}
	expectedEvent := integration_test_builder.NewEventBuilder(t).WithName("UserAFollowedUserBEvent").WithData(expectedUserAFollowedUserBEvent).Build()
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`

	controller.FollowUser(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
	integration_test_assert.AssertRelationshipExists(t, db, newUserPair)
}

func serializeData(data any) ([]byte, error) {
	return json.Marshal(data)
}
