package check_relationship

import (
	database "followservice/internal/db"
	model "followservice/internal/model/domain"
)

//go:generate mockgen -source=repository.go -destination=test/mock/repository.go

type CheckRelationshipRepository struct {
	dataRepository *database.Database
}

func NewCheckRelationshipRepository(dataRepository *database.Database) *CheckRelationshipRepository {
	return &CheckRelationshipRepository{
		dataRepository: dataRepository,
	}
}

func (r *CheckRelationshipRepository) CheckRelationship(userPair *model.UserPairRelationship) (bool, error) {
	return r.dataRepository.Client.RelationshipExists(userPair)
}
