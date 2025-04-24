package unit_test_unfollow_user

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"followservice/internal/bus"
	mock_bus "followservice/internal/bus/test/mock"
	"followservice/internal/features/unfollow_user"
	mock_unfollow_user "followservice/internal/features/unfollow_user/test/mock"
	model "followservice/internal/model/domain"
	"followservice/internal/model/events"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var serviceLoggerOutput bytes.Buffer
var serviceRepository *mock_unfollow_user.MockRepository
var serviceExternalBus *mock_bus.MockExternalBus
var serviceBus *bus.EventBus
var unfollowUserService *unfollow_user.UnfollowUserService

func setUpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	serviceRepository = mock_unfollow_user.NewMockRepository(ctrl)
	log.Logger = log.Output(&serviceLoggerOutput)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus = bus.NewEventBus(serviceExternalBus)
	unfollowUserService = unfollow_user.NewUnfollowUserService(serviceRepository, serviceBus)
}

func TestUnfollowUserWithService_WhenItReturnsSuccess(t *testing.T) {
	setUpService(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	expectedUserAUnfollowedUserBEvent := &events.UserAUnfollowedUserBEvent{
		FollowerID: newUserPair.FollowerID,
		FolloweeID: newUserPair.FolloweeID,
	}
	expectedEvent, _ := createEvent(events.UserAUnfollowedUserBEventName, expectedUserAUnfollowedUserBEvent)
	serviceRepository.EXPECT().RemoveUserRelationship(newUserPair).Return(nil)
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(nil)

	err := unfollowUserService.UnfollowUser(newUserPair)

	assert.Nil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "User pair relation was removed, "+newUserPair.FollowerID+" -> "+newUserPair.FolloweeID)
}

func TestErrorOnUnfollowUserWithService_WhenRemovingFromRepositoryFails(t *testing.T) {
	setUpService(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	serviceRepository.EXPECT().RemoveUserRelationship(newUserPair).Return(errors.New("some error"))

	err := unfollowUserService.UnfollowUser(newUserPair)

	assert.NotNil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Error removing user pair relation, "+newUserPair.FollowerID+" -> "+newUserPair.FolloweeID)
}

func TestErrorOnUnfollowUserWithService_WhenPublishingEventFails(t *testing.T) {
	setUpService(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	expectedUserAUnfollowedUserBEvent := &events.UserAUnfollowedUserBEvent{
		FollowerID: newUserPair.FollowerID,
		FolloweeID: newUserPair.FolloweeID,
	}
	expectedEvent, _ := createEvent(events.UserAUnfollowedUserBEventName, expectedUserAUnfollowedUserBEvent)
	serviceRepository.EXPECT().RemoveUserRelationship(newUserPair).Return(nil)
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(errors.New("some error"))

	err := unfollowUserService.UnfollowUser(newUserPair)

	assert.NotNil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Publishing UserAUnfollowedUserBEvent failed, "+expectedUserAUnfollowedUserBEvent.FollowerID+" -> "+expectedUserAUnfollowedUserBEvent.FolloweeID)
}

func createEvent(eventName string, eventData any) (*bus.Event, error) {
	dataEvent, err := serialize(eventData)
	if err != nil {
		return nil, err
	}

	return &bus.Event{
		Type: eventName,
		Data: dataEvent,
	}, nil
}

func serialize(data any) ([]byte, error) {
	return json.Marshal(data)
}
