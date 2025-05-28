package integration_test_check_relationship

import (
	database "followservice/internal/db"
	"followservice/internal/features/check_relationship"
	model "followservice/internal/model/domain"
	integration_test_arrange "followservice/test/integration_test_common/arrange"
	integration_test_assert "followservice/test/integration_test_common/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

var db *database.Database
var controller *check_relationship.CheckRelationshipController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUp(t *testing.T) {
	// Mocks
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)

	// Real infrastructure and services
	db = integration_test_arrange.CreateTestDatabase(t, ginContext)
	repository := check_relationship.NewCheckRelationshipRepository(db)
	service := check_relationship.NewCheckRelationshipService(repository)
	controller = check_relationship.NewCheckRelationshipController(service)
}

func tearDown() {
	db.Client.Clean()
}

func TestCheckRelationship_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp(t)
	defer tearDown()
	expectedUserPair := &model.UserPairRelationship{
		FollowerID: "user1",
		FolloweeID: "user2",
	}
	populateDb(t, expectedUserPair.FollowerID, expectedUserPair.FolloweeID)
	ginContext.Request, _ = http.NewRequest("GET", "/relationship/exists", nil)
	u := url.Values{}
	u.Add("followerId", expectedUserPair.FollowerID)
	u.Add("followeeId", expectedUserPair.FolloweeID)
	ginContext.Request.URL.RawQuery = u.Encode()
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": true
	}`

	controller.CheckRelationship(ginContext)

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
