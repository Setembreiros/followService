package unit_test_check_relationship

import (
	"errors"
	"testing"

	"followservice/internal/features/check_relationship"
	mock_check_relationship "followservice/internal/features/check_relationship/test/mock"
	model "followservice/internal/model/domain"

	"github.com/stretchr/testify/assert"
)

var serviceRepository *mock_check_relationship.MockRepository
var service check_relationship.Service

func setUpService(t *testing.T) {
	setUp(t)
	serviceRepository = mock_check_relationship.NewMockRepository(ctrl)
	service = check_relationship.NewCheckRelationshipService(serviceRepository)
}

func TestCheckRelationshipService_WhenSuccess(t *testing.T) {
	setUpService(t)
	userPair := &model.UserPairRelationship{FollowerID: "userA", FolloweeID: "userB"}
	serviceRepository.EXPECT().CheckRelationship(userPair).Return(true, nil)
	result, err := service.CheckRelationship(userPair)
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestCheckRelationshipService_WhenNotExists(t *testing.T) {
	setUpService(t)
	userPair := &model.UserPairRelationship{FollowerID: "userA", FolloweeID: "userB"}
	serviceRepository.EXPECT().CheckRelationship(userPair).Return(false, nil)
	result, err := service.CheckRelationship(userPair)
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestCheckRelationshipService_WhenError(t *testing.T) {
	setUpService(t)
	userPair := &model.UserPairRelationship{FollowerID: "userA", FolloweeID: "userB"}
	serviceRepository.EXPECT().CheckRelationship(userPair).Return(false, errors.New("db error"))
	_, err := service.CheckRelationship(userPair)
	assert.Error(t, err)
	assert.Contains(t, loggerOutput.String(), "Error checking relationship, userA -> userB")
}
