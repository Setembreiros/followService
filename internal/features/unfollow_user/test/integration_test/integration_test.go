package integration_test_unfollow_user

import (
	"encoding/json"
	"followservice/cmd/provider"
	"followservice/internal/bus"
	mock_bus "followservice/internal/bus/test/mock"
	database "followservice/internal/db"
	"followservice/internal/features/unfollow_user"
	model "followservice/internal/model/domain"
	"followservice/internal/model/events"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var db *database.Database
var controller *unfollow_user.UnfollowUserController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context
var serviceExternalBus *mock_bus.MockExternalBus

func setUpHandler(t *testing.T) {
	// Mocks
	ctrl := gomock.NewController(t)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus := bus.NewEventBus(serviceExternalBus)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)

	// Real infrastructure and services
	provider := provider.NewProvider("test")
	db = provider.ProvideDb(ginContext)
	repository := unfollow_user.NewUnfollowUserRepository(db)
	service := unfollow_user.NewUnfollowUserService(repository, serviceBus)
	controller = unfollow_user.NewUnfollowUserController(service)
}

func tearDown() {
	db.Client.Clean()
}

func TestUnfollowUser_WhenItReturnsSuccess(t *testing.T) {
	setUpHandler(t)
	defer tearDown()
	req, _ := http.NewRequest("DELETE", "/follow?followerId=usernameA&followeeId=usernameB", nil)
	ginContext.Request = req
	existingUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	AddRelationshipToDatabase(t, existingUserPair)
	expectedUserAUnfollowedUserBEvent := &events.UserAUnfollowedUserBEvent{
		FollowerID: existingUserPair.FollowerID,
		FolloweeID: existingUserPair.FolloweeID,
	}
	expectedEvent, _ := createEvent("UserAUnfollowedUserBEvent", expectedUserAUnfollowedUserBEvent)
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(nil)

	controller.UnfollowUser(ginContext)

	assertSuccessResult(t)
	assertRelationshipDoesNotExists(t, existingUserPair)
}

func AddRelationshipToDatabase(t *testing.T, userPair *model.UserPairRelationship) {
	err := db.Client.CreateRelationship(userPair)
	assert.Nil(t, err)
	existsInDatabase, err := db.Client.RelationshipExists(userPair)
	assert.Nil(t, err)
	assert.Equal(t, existsInDatabase, true)
}

func assertSuccessResult(t *testing.T) {
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": null
	}`
	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func assertRelationshipDoesNotExists(t *testing.T, userPair *model.UserPairRelationship) {
	existsInDatabase, err := db.Client.RelationshipExists(userPair)
	assert.Nil(t, err)
	assert.Equal(t, existsInDatabase, false)
}

func createEvent(eventName string, eventData any) (*bus.Event, error) {
	dataEvent, err := serializeData(eventData)
	if err != nil {
		return nil, err
	}

	return &bus.Event{
		Type: eventName,
		Data: dataEvent,
	}, nil
}

func serializeData(data any) ([]byte, error) {
	return json.Marshal(data)
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
