package integration_test_get_user_followers

import (
	"fmt"
	"followservice/cmd/provider"
	database "followservice/internal/db"
	"followservice/internal/features/get_user_followers"
	model "followservice/internal/model/domain"
	integration_test_assert "followservice/test/integration_test_common/assert"
	integration_test_builder "followservice/test/integration_test_common/builder"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

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
	db = getPopulatedDb(t)
	provider := provider.NewProvider("test")
	repository := get_user_followers.NewGetUserFollowersRepository(db, provider.ProvideCache(ginContext))
	service := get_user_followers.NewGetUserFollowersService(repository)
	controller = get_user_followers.NewGetUserFollowersController(service)
}

func tearDown() {
	db.Client.Clean()
}

func TestGetUserFollowers_WhenItReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	username := "usernameA"
	LastFollowerId := "usernameB"
	limit := 4
	req, _ := http.NewRequest("GET", fmt.Sprintf("/followers?username=%s&lastFollowerId=%s&limit=%d", username, LastFollowerId, limit), nil)
	ginContext.Request = req
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"followers":["usernameC","usernameD","usernameE","usernameF"],
			"lastFollowerId":"usernameF"
		}
	}`

	controller.GetUserFollowers(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
}

func getPopulatedDb(t *testing.T) *database.Database {
	existingUserPairs := []*model.UserPairRelationship{
		{
			FollowerID: "usernameB",
			FolloweeID: "usernameA",
		},
		{
			FollowerID: "usernameC",
			FolloweeID: "usernameA",
		},
		{
			FollowerID: "usernameD",
			FolloweeID: "usernameA",
		},
		{
			FollowerID: "usernameE",
			FolloweeID: "usernameA",
		},
		{
			FollowerID: "usernameF",
			FolloweeID: "usernameA",
		},
		{
			FollowerID: "usernameG",
			FolloweeID: "usernameA",
		},
	}

	dbBuilder := integration_test_builder.NewDatabaseBuilder(t, ginContext)
	for _, existingUserPair := range existingUserPairs {
		dbBuilder.WithRelationship(existingUserPair)
	}

	return dbBuilder.Build()
}
