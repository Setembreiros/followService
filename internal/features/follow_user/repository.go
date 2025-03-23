package follow_user

import (
	database "followservice/internal/db"
	model "followservice/internal/model/domain"
)

type FollowUserRepository struct {
	dataRepository *database.Database
}

func NewFollowUserRepository(dataRepository *database.Database) *FollowUserRepository {
	return &FollowUserRepository{
		dataRepository: dataRepository,
	}
}

func (r *FollowUserRepository) AddUserRelationship(data *model.UserPairRelationship) error {
	return r.dataRepository.Client.CreateRelationship(data)
}
