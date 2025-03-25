package integration_test_get_user_followers

import (
	"fmt"
	database "followservice/internal/db"
	"followservice/internal/features/get_user_followers"
	model "followservice/internal/model/domain"
	integration_test_arrange "followservice/test/integration_test_common/arrange"
	integration_test_assert "followservice/test/integration_test_common/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

var cache *database.Cache
var db *database.Database
var controller *get_user_followers.GetUserFollowersController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUp(t *testing.T) {
	// Mocks
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)

	// Real infrastructure and services
	cache = integration_test_arrange.CreateTestCache(t, ginContext)
	db = integration_test_arrange.CreateTestDatabase(t, ginContext)
	repository := get_user_followers.NewGetUserFollowersRepository(db, cache)
	service := get_user_followers.NewGetUserFollowersService(repository)
	controller = get_user_followers.NewGetUserFollowersController(service)
}

func tearDown() {
	db.Client.Clean()
	cache.Client.Clean()
}

func TestGetUserFollowers_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	username := "username1"
	lastFollowerId := "username2"
	limit := 4
	populateDb(t, username, lastFollowerId)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/followers?username=%s&lastFollowerId=%s&limit=%d", username, lastFollowerId, limit), nil)
	ginContext.Request = req
	expectedFollowers := []string{"usernameA", "usernameB", "usernameC", "usernameD"}
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"followers":["usernameA", "usernameB", "usernameC", "usernameD"],
			"lastFollowerId":"usernameD"
		}
	}`

	controller.GetUserFollowers(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
	integration_test_assert.AssertCachedUserFollowersExists(t, cache, username, lastFollowerId, limit, expectedFollowers)
}

func TestGetUserFollowers_WhenCacheReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	username := "username1"
	LastFollowerId := "username2"
	limit := 4
	populateCache(t, username, LastFollowerId, limit)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/followers?username=%s&lastFollowerId=%s&limit=%d", username, LastFollowerId, limit), nil)
	ginContext.Request = req
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"followers":["usernameA","usernameB","usernameC","usernameD"],
			"lastFollowerId":"usernameD"
		}
	}`

	controller.GetUserFollowers(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
}

func populateDb(t *testing.T, followeeId, lastFollowerId string) {
	existingUserPairs := []*model.UserPairRelationship{
		{
			FollowerID: lastFollowerId,
			FolloweeID: followeeId,
		},
		{
			FollowerID: "usernameA",
			FolloweeID: followeeId,
		},
		{
			FollowerID: "usernameB",
			FolloweeID: followeeId,
		},
		{
			FollowerID: "usernameC",
			FolloweeID: followeeId,
		},
		{
			FollowerID: "usernameD",
			FolloweeID: followeeId,
		},
		{
			FollowerID: "usernameE",
			FolloweeID: followeeId,
		},
	}

	for _, existingUserPair := range existingUserPairs {
		integration_test_arrange.AddRelationshipToDatabase(t, db, existingUserPair)
	}
}

func populateCache(t *testing.T, followeeId, lastFollowerId string, limit int) {
	integration_test_arrange.AddCachedFollowersToCache(t, cache, followeeId, lastFollowerId, limit, []string{"usernameA", "usernameB", "usernameC", "usernameD"})
}
