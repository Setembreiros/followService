package integration_test_follow_user

import (
	"bytes"
	"encoding/json"
	"followservice/cmd/provider"
	"followservice/internal/bus"
	mock_bus "followservice/internal/bus/test/mock"
	database "followservice/internal/db"
	"followservice/internal/features/follow_user"
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
var controller *follow_user.FollowUserController
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
	repository := follow_user.NewFollowUserRepository(db)
	service := follow_user.NewFollowUserService(repository, serviceBus)
	controller = follow_user.NewFollowUserController(service)
}

func tearDown() {
	db.Client.Clean()
}

func TestFollowUser_WhenItReturnsSuccess(t *testing.T) {
	setUpHandler(t)
	defer tearDown()
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	expectedUserAFollowedUserBEvent := &events.UserAFollowedUserBEvent{
		FollowerID: newUserPair.FollowerID,
		FolloweeID: newUserPair.FolloweeID,
	}
	expectedEvent, _ := createEvent("UserAFollowedUserBEvent", expectedUserAFollowedUserBEvent)
	data, _ := serializeData(newUserPair)
	ginContext.Request = httptest.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(data))
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(nil)

	controller.FollowUser(ginContext)

	assertSuccessResult(t)
	assertRelationshipExists(t, newUserPair)
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

func assertRelationshipExists(t *testing.T, userPair *model.UserPairRelationship) {
	existsInDatabase, err := db.Client.RelationshipExists(userPair)
	assert.Nil(t, err)
	assert.Equal(t, existsInDatabase, true)
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
