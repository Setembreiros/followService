package follow_user_artist

import (
	"followservice/internal/bus"
	model "followservice/internal/model/domain"
	"followservice/internal/model/events"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	AddUserRelation(data *model.UserPairRelationship) error
}

type FollowUserArtistService struct {
	repository Repository
	bus        *bus.EventBus
}

func NewFollowUserArtistService(repository Repository, bus *bus.EventBus) *FollowUserArtistService {
	return &FollowUserArtistService{
		repository: repository,
		bus:        bus,
	}
}

func (s *FollowUserArtistService) FollowUserArtist(userPair *model.UserPairRelationship) error {
	err := s.repository.AddUserRelation(userPair)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error adding user pair relation, %s -> %s", userPair.FollowerID, userPair.FolloweeID)
		return err
	}

	err = s.publishPostWasCreatedEvent(userPair)
	if err != nil {
		return err
	}

	log.Error().Stack().Err(err).Msgf("User pair relation was created, %s -> %s", userPair.FollowerID, userPair.FolloweeID)

	return nil
}

func (s *FollowUserArtistService) publishPostWasCreatedEvent(data *model.UserPairRelationship) error {
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
