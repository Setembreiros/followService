package unfollow_user

import (
	"followservice/internal/api"
	model "followservice/internal/model/domain"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=controller.go -destination=test/mock/controller.go

type UnfollowUserController struct {
	service Service
}

type Service interface {
	UnfollowUser(userPair *model.UserPairRelationship) error
}

func NewUnfollowUserController(service Service) *UnfollowUserController {
	return &UnfollowUserController{
		service: service,
	}
}

func (controller *UnfollowUserController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.DELETE("/follow", controller.UnfollowUser)
}

func (controller *UnfollowUserController) UnfollowUser(c *gin.Context) {
	log.Info().Msg("Handling Request DELETE UnfollowUser")
	followerId := c.Query("followerId")
	if followerId == "" {
		api.SendBadRequest(c, "Missing followerId parameter")
		return
	}
	followeeId := c.Query("followeeId")
	if followeeId == "" {
		api.SendBadRequest(c, "Missing followeeId parameter")
		return
	}
	userPair := &model.UserPairRelationship{
		FollowerID: followerId,
		FolloweeID: followeeId,
	}

	err := controller.service.UnfollowUser(userPair)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOK(c)
}
