package get_user_followers

import (
	database "followservice/internal/db"
)

type GetUserFollowersRepository struct {
	dataRepository *database.Database
}

func NewGetUserFollowersRepository(dataRepository *database.Database) *GetUserFollowersRepository {
	return &GetUserFollowersRepository{
		dataRepository: dataRepository,
	}
}

func (r *GetUserFollowersRepository) GetUserFollowers(username string, lastFollowerId string, limit int) ([]string, string, error) {
	followers, newLastFollowerId, err := r.dataRepository.Client.GetUserFollowers(username, lastFollowerId, limit)
	if err != nil {
		return []string{}, "", err
	}

	return followers, newLastFollowerId, nil
}
