package follow_user_artist

import (
	database "followservice/internal/db"
	model "followservice/internal/model/domain"
)

type FollowUserArtistRepository struct {
	dataRepository *database.Database
}

func NewFollowUserArtistRepository(dataRepository *database.Database) *FollowUserArtistRepository {
	return &FollowUserArtistRepository{
		dataRepository: dataRepository,
	}
}

func (r *FollowUserArtistRepository) AddUserRelationship(data *model.UserPairRelationship) error {
	return r.dataRepository.Client.CreateRelationship(data)
}
