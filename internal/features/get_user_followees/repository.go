package get_user_followees

import (
	database "followservice/internal/db"
)

type GetUserFolloweesRepository struct {
	dataRepository *database.Database
}

func NewGetUserFolloweesRepository(dataRepository *database.Database) *GetUserFolloweesRepository {
	return &GetUserFolloweesRepository{
		dataRepository: dataRepository,
	}
}

func (r *GetUserFolloweesRepository) GetUserFollowees(username string, lastFolloweeId string, limit int) ([]string, string, error) {
	followees, newLastFolloweeId, err := r.dataRepository.Client.GetUserFollowees(username, lastFolloweeId, limit)
	if err != nil {
		return []string{}, "", err
	}

	return followees, newLastFolloweeId, nil
}
