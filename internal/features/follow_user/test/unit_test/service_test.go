package unit_test_follow_user

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"followservice/internal/bus"
	mock_bus "followservice/internal/bus/test/mock"
	"followservice/internal/features/follow_user"
	mock_follow_user "followservice/internal/features/follow_user/test/mock"
	model "followservice/internal/model/domain"
	"followservice/internal/model/events"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var serviceLoggerOutput bytes.Buffer
var serviceRepository *mock_follow_user.MockRepository
var serviceExternalBus *mock_bus.MockExternalBus
var serviceBus *bus.EventBus
var followUserService *follow_user.FollowUserService

func setUpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	serviceRepository = mock_follow_user.NewMockRepository(ctrl)
	log.Logger = log.Output(&serviceLoggerOutput)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus = bus.NewEventBus(serviceExternalBus)
	followUserService = follow_user.NewFollowUserService(serviceRepository, serviceBus)
}

func TestFollowUserWithService_WhenItReturnsSuccess(t *testing.T) {
	setUpService(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	expectedUserAFollowedUserBEvent := &events.UserAFollowedUserBEvent{
		FollowerID: newUserPair.FollowerID,
		FolloweeID: newUserPair.FolloweeID,
	}
	expectedEvent, _ := createEvent("UserAFollowedUserBEvent", expectedUserAFollowedUserBEvent)
	serviceRepository.EXPECT().AddUserRelationship(newUserPair).Return(nil)
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(nil)

	err := followUserService.FollowUser(newUserPair)

	assert.Nil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "User pair relation was created, "+newUserPair.FollowerID+" -> "+newUserPair.FolloweeID)
}

func TestErrorOnFollowUserWithService_WhenAddingToRepositoryFails(t *testing.T) {
	setUpService(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	serviceRepository.EXPECT().AddUserRelationship(newUserPair).Return(errors.New("some error"))

	err := followUserService.FollowUser(newUserPair)

	assert.NotNil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Error adding user pair relation, "+newUserPair.FollowerID+" -> "+newUserPair.FolloweeID)
}

func TestErrorOnFollowUserWithService_WhenPublishingEventFails(t *testing.T) {
	setUpService(t)
	newUserPair := &model.UserPairRelationship{
		FollowerID: "usernameA",
		FolloweeID: "usernameB",
	}
	expectedUserAFollowedUserBEvent := &events.UserAFollowedUserBEvent{
		FollowerID: newUserPair.FollowerID,
		FolloweeID: newUserPair.FolloweeID,
	}
	expectedEvent, _ := createEvent("UserAFollowedUserBEvent", expectedUserAFollowedUserBEvent)
	serviceRepository.EXPECT().AddUserRelationship(newUserPair).Return(nil)
	serviceExternalBus.EXPECT().Publish(expectedEvent).Return(errors.New("some error"))

	err := followUserService.FollowUser(newUserPair)

	assert.NotNil(t, err)
	assert.Contains(t, serviceLoggerOutput.String(), "Publishing UserAFollowedUserBEvent failed, "+expectedUserAFollowedUserBEvent.FollowerID+" -> "+expectedUserAFollowedUserBEvent.FolloweeID)
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
