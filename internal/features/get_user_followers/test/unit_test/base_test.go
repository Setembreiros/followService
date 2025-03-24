package unit_test_get_user_followers

import (
	"bytes"
	"followservice/internal/features/get_user_followers"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

var ctrl *gomock.Controller
var loggerOutput bytes.Buffer
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUp(t *testing.T) {
	ctrl = gomock.NewController(t)
	log.Logger = log.Output(&loggerOutput)
	controller = get_user_followers.NewGetUserFollowersController(controllerService)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
