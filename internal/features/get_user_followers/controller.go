package get_user_followers

import (
	"followservice/internal/api"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=controller.go -destination=test/mock/controller.go

type GetUserFollowersController struct {
	service Service
}

type Service interface {
	GetUserFollowers(username string, lastPostId string, limit int) ([]string, string, error)
}

type GetUserFollowersResponse struct {
	Followers      []string `json:"followers"`
	LastFollowerId string   `json:"lastFollowerId"`
}

func NewGetUserFollowersController(service Service) *GetUserFollowersController {
	return &GetUserFollowersController{
		service: service,
	}
}

func (controller *GetUserFollowersController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/followers", controller.GetUserFollowers)
}

func (controller *GetUserFollowersController) GetUserFollowers(c *gin.Context) {
	log.Info().Msg("Handling Request GET GetUserFollowers")
	username, lastFollowerId, limit := getQueryParameters(c)
	if username == "" {
		return
	}

	followers, lastFollowerId, err := controller.service.GetUserFollowers(username, lastFollowerId, limit)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, &GetUserFollowersResponse{
		Followers:      followers,
		LastFollowerId: lastFollowerId,
	})
}

func getQueryParameters(c *gin.Context) (string, string, int) {
	username := c.Query("username")
	if username == "" {
		api.SendBadRequest(c, "Missing username parameter")
		return "", "", 0
	}
	lastFollowerId := c.DefaultQuery("lastFollowerId", "")
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "12"))
	if err != nil || limit <= 0 {
		api.SendBadRequest(c, "Invalid pagination parameters, limit has to be greater than 0")
		return "", "", 0
	}

	return username, lastFollowerId, limit
}
