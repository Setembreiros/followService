package integration_test_assert

import (
	database "followservice/internal/db"
	model "followservice/internal/model/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertRelationshipDoesNotExists(t *testing.T, db *database.Database, expectedUserPair *model.UserPairRelationship) {
	existsInDatabase, err := db.Client.RelationshipExists(expectedUserPair)
	assert.Nil(t, err)
	assert.Equal(t, existsInDatabase, false)
}
