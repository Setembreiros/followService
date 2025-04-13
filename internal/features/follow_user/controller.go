package follow_user

import (
	"followservice/internal/api"
	model "followservice/internal/model/domain"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=controller.go -destination=test/mock/controller.go

type FollowUserController struct {
	service Service
}

type Service interface {
	FollowUser(userPair *model.UserPairRelationship) error
}

func NewFollowUserController(service Service) *FollowUserController {
	return &FollowUserController{
		service: service,
	}
}

func (controller *FollowUserController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/follow", controller.FollowUser)
}

func (controller *FollowUserController) FollowUser(c *gin.Context) {
	log.Info().Msg("Handling Request POST FollowUser")
	var userPair model.UserPairRelationship

	if err := c.BindJSON(&userPair); err != nil {
		log.Error().Stack().Err(err).Msg("Invalid Data")
		api.SendBadRequest(c, "Invalid Json Request")
		return
	}

	err := controller.service.FollowUser(&userPair)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOK(c)
}
