package unfollow_user

import (
	database "followservice/internal/db"
	model "followservice/internal/model/domain"
)

type UnfollowUserRepository struct {
	dataRepository *database.Database
}

func NewUnfollowUserRepository(dataRepository *database.Database) *UnfollowUserRepository {
	return &UnfollowUserRepository{
		dataRepository: dataRepository,
	}
}

func (r *UnfollowUserRepository) RemoveUserRelationship(data *model.UserPairRelationship) error {
	return r.dataRepository.Client.DeleteRelationship(data)
}
