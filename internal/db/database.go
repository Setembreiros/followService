package database

import model "followservice/internal/model/domain"

//go:generate mockgen -source=database.go -destination=test/mock/database.go

type Database struct {
	Client DatabaseClient
}

type DatabaseClient interface {
	Clean()
	RelationshipExists(userPair *model.UserPairRelationship) (bool, error)
	CreateRelationship(relationship *model.UserPairRelationship) error
	DeleteRelationship(relationship *model.UserPairRelationship) error
}

func NewDatabase(client DatabaseClient) *Database {
	return &Database{
		Client: client,
	}
}
