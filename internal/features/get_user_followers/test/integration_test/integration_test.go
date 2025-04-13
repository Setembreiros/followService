package integration_test_get_user_followers

import (
	database "followservice/internal/db"
	"followservice/internal/features/get_user_followers"
	model "followservice/internal/model/domain"
	integration_test_arrange "followservice/test/integration_test_common/arrange"
	integration_test_assert "followservice/test/integration_test_common/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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
	username := "USERC"
	lastFollowerId := "USER"
	limit := 4
	populateDb(t, username, lastFollowerId)
	ginContext.Request, _ = http.NewRequest("GET", "/followers", nil)
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	u := url.Values{}
	u.Add("lastFollowerId", lastFollowerId)
	u.Add("limit", strconv.Itoa(limit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedFollowers := []string{"USERA", "USERB"}
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"followers":["USERA", "USERB"],
			"lastFollowerId":"USERB"
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
	lastFollowerId := "username2"
	limit := 4
	populateCache(t, username, lastFollowerId, limit)
	ginContext.Request, _ = http.NewRequest("GET", "/followers", nil)
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	u := url.Values{}
	u.Add("lastFollowerId", lastFollowerId)
	u.Add("limit", strconv.Itoa(limit))
	ginContext.Request.URL.RawQuery = u.Encode()
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
			FollowerID: "USERA",
			FolloweeID: followeeId,
		},
		{
			FollowerID: "USERB",
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
