package integration_test_builder

import (
	"context"
	"followservice/cmd/provider"
	database "followservice/internal/db"
	model "followservice/internal/model/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DatabaseBuilder struct {
	db            *database.Database
	t             *testing.T
	relationships []*model.UserPairRelationship
}

func NewDatabaseBuilder(t *testing.T, ctx context.Context) *DatabaseBuilder {
	provider := provider.NewProvider("test")
	return &DatabaseBuilder{
		db:            provider.ProvideDb(ctx),
		t:             t,
		relationships: []*model.UserPairRelationship{},
	}
}

func (d *DatabaseBuilder) WithRelationship(relationship *model.UserPairRelationship) *DatabaseBuilder {
	d.relationships = append(d.relationships, relationship)
	return d
}

func (d *DatabaseBuilder) Build() *database.Database {
	for _, relationship := range d.relationships {
		d.addRelationshipToDatabase(relationship)
	}
	return d.db
}

func (d *DatabaseBuilder) addRelationshipToDatabase(userPair *model.UserPairRelationship) {
	err := d.db.Client.CreateRelationship(userPair)
	assert.Nil(d.t, err)
	existsInDatabase, err := d.db.Client.RelationshipExists(userPair)
	assert.Nil(d.t, err)
	assert.Equal(d.t, existsInDatabase, true)
}
