package get_user_followers

import (
	database "followservice/internal/db"
)

type GetUserFollowersRepository struct {
	dataRepository *database.Database
	cache          *database.Cache
}

func NewGetUserFollowersRepository(dataRepository *database.Database, cache *database.Cache) *GetUserFollowersRepository {
	return &GetUserFollowersRepository{
		dataRepository: dataRepository,
		cache:          cache,
	}
}

func (r *GetUserFollowersRepository) GetUserFollowers(username string, lastFollowerId string, limit int) ([]string, string, error) {
	followers, newLastFollowerId, found := r.cache.Client.GetUserFollowers(username, lastFollowerId, limit)
	if found {
		return followers, lastFollowerId, nil
	}

	followers, newLastFollowerId, err := r.dataRepository.Client.GetUserFollowers(username, lastFollowerId, limit)
	if err != nil {
		return []string{}, "", err
	}

	return followers, newLastFollowerId, nil
}
