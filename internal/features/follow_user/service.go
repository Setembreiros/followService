package follow_user

import (
	"followservice/internal/bus"
	model "followservice/internal/model/domain"
	"followservice/internal/model/events"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	AddUserRelationship(data *model.UserPairRelationship) error
}

type FollowUserService struct {
	repository Repository
	bus        *bus.EventBus
}

func NewFollowUserService(repository Repository, bus *bus.EventBus) *FollowUserService {
	return &FollowUserService{
		repository: repository,
		bus:        bus,
	}
}

func (s *FollowUserService) FollowUser(userPair *model.UserPairRelationship) error {
	err := s.repository.AddUserRelationship(userPair)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error adding user pair relation, %s -> %s", userPair.FollowerID, userPair.FolloweeID)
		return err
	}

	err = s.publishUserAFollowedUserBEvent(userPair)
	if err != nil {
		return err
	}

	log.Info().Msgf("User pair relation was created, %s -> %s", userPair.FollowerID, userPair.FolloweeID)

	return nil
}

func (s *FollowUserService) publishUserAFollowedUserBEvent(data *model.UserPairRelationship) error {
	event := &events.UserAFollowedUserBEvent{
		FollowerID: data.FollowerID,
		FolloweeID: data.FolloweeID,
	}

	err := s.bus.Publish("UserAFollowedUserBEvent", event)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Publishing UserAFollowedUserBEvent failed, %s -> %s", event.FollowerID, event.FolloweeID)
		return err
	}

	return nil
}
