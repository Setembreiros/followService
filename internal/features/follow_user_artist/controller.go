package follow_user_artist

import (
	"followservice/internal/api"
	model "followservice/internal/model/domain"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=controller.go -destination=test/mock/controller.go

type FollowUserArtistController struct {
	service Service
}

type Service interface {
	FollowUserArtist(userPair *model.UserPairRelationship) error
}

func NewFollowUserArtistController(service Service) *FollowUserArtistController {
	return &FollowUserArtistController{
		service: service,
	}
}

func (controller *FollowUserArtistController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/post", controller.FollowUserArtist)
}

func (controller *FollowUserArtistController) FollowUserArtist(c *gin.Context) {
	log.Info().Msg("Handling Request POST FollowUserArtist")
	var userPair model.UserPairRelationship

	if err := c.BindJSON(&userPair); err != nil {
		log.Error().Stack().Err(err).Msg("Invalid Data")
		api.SendBadRequest(c, "Invalid Json Request")
		return
	}

	err := controller.service.FollowUserArtist(&userPair)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOK(c)
}
