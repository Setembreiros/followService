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
	relationshipExists, err := r.dataRepository.Client.RelationshipExists(data)
	if err != nil {
		return err
	}
	if relationshipExists {
		return database.NewRelationshipAlreadyExistsError(data.FollowerID, data.FolloweeID)
	}
	return r.dataRepository.Client.CreateRelationship(data)
}
