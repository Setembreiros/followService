package integration_test_unfollow_user

import (
	"followservice/internal/bus"
	mock_bus "followservice/internal/bus/test/mock"
	database "followservice/internal/db"
	"followservice/internal/features/unfollow_user"
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
var controller *unfollow_user.UnfollowUserController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context
var serviceExternalBus *mock_bus.MockExternalBus
var existingUserPair *model.UserPairRelationship

func setUp(t *testing.T) {
	// Mocks
	ctrl := gomock.NewController(t)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus := bus.NewEventBus(serviceExternalBus)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)

	// Real infrastructure and services
	existingUserPair = &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	db = integration_test_builder.NewDatabaseBuilder(t, ginContext).WithRelationship(existingUserPair).Build()
	repository := unfollow_user.NewUnfollowUserRepository(db)
	service := unfollow_user.NewUnfollowUserService(repository, serviceBus)
	controller = unfollow_user.NewUnfollowUserController(service)
}

func tearDown() {
	db.Client.Clean()
}

func TestUnfollowUser_WhenItReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	req, _ := http.NewRequest("DELETE", "/follow?followerId=usernameA&followeeId=usernameB", nil)
	ginContext.Request = req
	expectedUserAUnfollowedUserBEvent := &events.UserAUnfollowedUserBEvent{
		FollowerID: existingUserPair.FollowerID,
		FolloweeID: existingUserPair.FolloweeID,
	}
	expectedEvent := integration_test_builder.NewEventBuilder(t).WithName("UserAUnfollowedUserBEvent").WithData(expectedUserAUnfollowedUserBEvent).Build()
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`

	controller.UnfollowUser(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
	integration_test_assert.AssertRelationshipDoesNotExists(t, db, existingUserPair)
}
