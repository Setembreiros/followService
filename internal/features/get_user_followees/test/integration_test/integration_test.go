package integration_test_get_user_followees

import (
	database "followservice/internal/db"
	"followservice/internal/features/get_user_followees"
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

var db *database.Database
var controller *get_user_followees.GetUserFolloweesController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUp(t *testing.T) {
	// Mocks
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)

	// Real infrastructure and services
	db = integration_test_arrange.CreateTestDatabase(t, ginContext)
	repository := get_user_followees.NewGetUserFolloweesRepository(db)
	service := get_user_followees.NewGetUserFolloweesService(repository)
	controller = get_user_followees.NewGetUserFolloweesController(service)
}

func tearDown() {
	db.Client.Clean()
}

func TestGetUserFollowees_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	username := "USERC"
	lastFolloweeId := "USER"
	limit := 4
	populateDb(t, username, lastFolloweeId)
	ginContext.Request, _ = http.NewRequest("GET", "/followees", nil)
	ginContext.Params = []gin.Param{{Key: "username", Value: username}}
	u := url.Values{}
	u.Add("lastFolloweeId", lastFolloweeId)
	u.Add("limit", strconv.Itoa(limit))
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"followees":["USERA", "USERB"],
			"lastFolloweeId":"USERB"
		}
	}`

	controller.GetUserFollowees(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
}

func populateDb(t *testing.T, followerId, lastFolloweeId string) {
	existingUserPairs := []*model.UserPairRelationship{
		{
			FollowerID: followerId,
			FolloweeID: lastFolloweeId,
		},
		{
			FollowerID: followerId,
			FolloweeID: "USERA",
		},
		{
			FollowerID: followerId,
			FolloweeID: "USERB",
		},
	}

	for _, existingUserPair := range existingUserPairs {
		integration_test_arrange.AddRelationshipToDatabase(t, db, existingUserPair)
	}
}
