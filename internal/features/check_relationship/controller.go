package check_relationship

import (
	"followservice/internal/api"
	model "followservice/internal/model/domain"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=controller.go -destination=test/mock/controller.go

type CheckRelationshipController struct {
	service Service
}

type Service interface {
	CheckRelationship(userPair *model.UserPairRelationship) (bool, error)
}

func NewCheckRelationshipController(service Service) *CheckRelationshipController {
	return &CheckRelationshipController{
		service: service,
	}
}

func (controller *CheckRelationshipController) Routes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/relationship/exists", controller.CheckRelationship)
}

func (controller *CheckRelationshipController) CheckRelationship(c *gin.Context) {
	log.Info().Msg("Handling Request Get CheckRelationship")

	followerId := c.Query("followerId")
	followeeId := c.Query("followeeId")

	if followerId == "" {
		api.SendBadRequest(c, "Missing followerId")
		return
	}

	if followeeId == "" {
		api.SendBadRequest(c, "Missing followeeId")
		return
	}

	userPair := &model.UserPairRelationship{
		FollowerID: followerId,
		FolloweeID: followeeId,
	}

	exists, err := controller.service.CheckRelationship(userPair)
	if err != nil {
		api.SendInternalServerError(c, err.Error())
		return
	}

	api.SendOKWithResult(c, exists)
}
