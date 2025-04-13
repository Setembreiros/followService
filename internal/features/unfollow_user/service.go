package unfollow_user

import (
	"followservice/internal/bus"
	model "followservice/internal/model/domain"
	"followservice/internal/model/events"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	RemoveUserRelationship(data *model.UserPairRelationship) error
}

type UnfollowUserService struct {
	repository Repository
	bus        *bus.EventBus
}

func NewUnfollowUserService(repository Repository, bus *bus.EventBus) *UnfollowUserService {
	return &UnfollowUserService{
		repository: repository,
		bus:        bus,
	}
}

func (s *UnfollowUserService) UnfollowUser(userPair *model.UserPairRelationship) error {
	err := s.repository.RemoveUserRelationship(userPair)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error removing user pair relation, %s -> %s", userPair.FollowerID, userPair.FolloweeID)
		return err
	}

	err = s.publishUserAUnfollowedUserBEvent(userPair)
	if err != nil {
		return err
	}

	log.Info().Msgf("User pair relation was removed, %s -> %s", userPair.FollowerID, userPair.FolloweeID)

	return nil
}

func (s *UnfollowUserService) publishUserAUnfollowedUserBEvent(data *model.UserPairRelationship) error {
	event := &events.UserAUnfollowedUserBEvent{
		FollowerID: data.FollowerID,
		FolloweeID: data.FolloweeID,
	}

	err := s.bus.Publish("UserAUnfollowedUserBEvent", event)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Publishing UserAUnfollowedUserBEvent failed, %s -> %s", event.FollowerID, event.FolloweeID)
		return err
	}

	return nil
}
