package unit_test_check_relationship

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"followservice/internal/features/check_relationship"
	mock_check_relationship "followservice/internal/features/check_relationship/test/mock"
	model "followservice/internal/model/domain"

	"github.com/stretchr/testify/assert"
)

var controllerService *mock_check_relationship.MockService
var controller *check_relationship.CheckRelationshipController

func setUpHandler(t *testing.T) {
	setUp(t)
	controllerService = mock_check_relationship.NewMockService(ctrl)
	controller = check_relationship.NewCheckRelationshipController(controllerService)
}

func TestCheckRelationshipWithController_WhenSuccess(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/relationship/exists", nil)
	expectedUserPair := &model.UserPairRelationship{
		FollowerID: "userA",
		FolloweeID: "userB",
	}
	u := url.Values{}
	u.Add("followerId", expectedUserPair.FollowerID)
	u.Add("followeeId", expectedUserPair.FolloweeID)
	ginContext.Request.URL.RawQuery = u.Encode()
	controllerService.EXPECT().CheckRelationship(expectedUserPair).Return(true, nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": true
	}`

	controller.CheckRelationship(ginContext)

	assert.Equal(t, 200, apiResponse.Code)
	assert.Equal(t, removeSpace(expectedBodyResponse), removeSpace(apiResponse.Body.String()))
}

func TestCheckRelationshipWithController_WhenMissingFollowerId(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/relationship/exists", nil)
	u := url.Values{}
	u.Add("followeeId", "userB")
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedBodyResponse := `{
		"error": true,
		"message": "Missing followerId",
		"content": null
	}`

	controller.CheckRelationship(ginContext)

	assert.Equal(t, 400, apiResponse.Code)
	assert.Equal(t, removeSpace(expectedBodyResponse), removeSpace(apiResponse.Body.String()))
}

func TestCheckRelationshipWithController_WhenMissingFolloweeId(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/relationship/exists", nil)
	u := url.Values{}
	u.Add("followerId", "userA")
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedBodyResponse := `{
		"error": true,
		"message": "Missing followeeId",
		"content": null
	}`

	controller.CheckRelationship(ginContext)

	assert.Equal(t, 400, apiResponse.Code)
	assert.Equal(t, removeSpace(expectedBodyResponse), removeSpace(apiResponse.Body.String()))
}

func TestCheckRelationshipWithController_WhenError(t *testing.T) {
	setUpHandler(t)
	ginContext.Request, _ = http.NewRequest("GET", "/relationship/exists", nil)
	expectedUserProfile := &model.UserPairRelationship{
		FollowerID: "userA",
		FolloweeID: "userB",
	}
	u := url.Values{}
	u.Add("followerId", expectedUserProfile.FollowerID)
	u.Add("followeeId", expectedUserProfile.FolloweeID)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedBodyResponse := `{
		"error": true,
		"message": "some error",
		"content": null
	}`
	controllerService.EXPECT().CheckRelationship(expectedUserProfile).Return(false, errors.New("some error"))

	controller.CheckRelationship(ginContext)

	assert.Equal(t, 500, apiResponse.Code)
	assert.Equal(t, removeSpace(expectedBodyResponse), removeSpace(apiResponse.Body.String()))
}
