package check_relationship

import (
	model "followservice/internal/model/domain"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	CheckRelationship(userPair *model.UserPairRelationship) (bool, error)
}

type CheckRelationshipService struct {
	repo Repository
}

func NewCheckRelationshipService(repo Repository) Service {
	return &CheckRelationshipService{repo: repo}
}

func (s *CheckRelationshipService) CheckRelationship(userPair *model.UserPairRelationship) (bool, error) {
	exists, err := s.repo.CheckRelationship(userPair)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error checking relationship, %s -> %s", userPair.FollowerID, userPair.FolloweeID)
		return false, err
	}

	return exists, nil
}
