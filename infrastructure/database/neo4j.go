package database

import (
	"context"

	database "followservice/internal/db"
	model "followservice/internal/model/domain"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
)

type Neo4jDBClient struct {
	ctx        context.Context
	dbUri      string
	dbUser     string
	dbPassword string
}

func NewNeo4jClient(dbUri, dbUser, dbPassword string, ctx context.Context) *Neo4jDBClient {
	client := &Neo4jDBClient{
		dbUri:      dbUri,
		dbUser:     dbUser,
		dbPassword: dbPassword,
		ctx:        ctx,
	}

	client.verifyConnection()

	return client
}

func (c Neo4jDBClient) verifyConnection() {
	driver, err := c.getNeo4jDriver()
	if err != nil {
		panic(err)
	}
	defer driver.Close(c.ctx)

	err = driver.VerifyConnectivity(c.ctx)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Conection to Neo4j not stablished")
		panic(err)
	}
	log.Info().Msg("Connection to Neo4j established.")
}

func (c *Neo4jDBClient) RelationshipExists(userPair *model.UserPairRelationship) (bool, error) {
	driver, err := c.getNeo4jDriver()
	if err != nil {
		return false, err
	}
	defer driver.Close(c.ctx)
	session := driver.NewSession(c.ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(c.ctx)

	result, err := session.Run(c.ctx,
		`MATCH (follower:User {id: $followerId})-[r:FOLLOWS]->(followee:User {id: $followeeId})
		 RETURN count(r) AS relationshipCount`,
		map[string]interface{}{
			"followerId": userPair.FollowerID,
			"followeeId": userPair.FolloweeID,
		},
	)
	if err != nil {
		return false, err
	}

	// Verificar que existe polo menos unha relación
	if result.Next(c.ctx) {
		count := result.Record().Values[0].(int64)
		if count >= 1 {
			return true, nil
		}
	}

	return false, nil
}

func (c *Neo4jDBClient) CreateRelationship(userPair *model.UserPairRelationship) error {
	driver, err := c.getNeo4jDriver()
	if err != nil {
		return err
	}
	defer driver.Close(c.ctx)
	session := driver.NewSession(c.ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(c.ctx)

	// Create the relationship with a transaction
	_, err = session.ExecuteWrite(c.ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		// Cypher query to create or merge users and establish FOLLOWS relationship
		query := `
            MERGE (follower:User {id: $followerId})
            MERGE (followee:User {id: $followeeId})
            MERGE (follower)-[r:FOLLOWS]->(followee)
            ON CREATE SET r.createdAt = datetime()
            RETURN r
        `
		params := map[string]interface{}{
			"followerId": userPair.FollowerID,
			"followeeId": userPair.FolloweeID,
		}

		result, err := tx.Run(c.ctx, query, params)
		if err != nil {
			return nil, err
		}

		// Check if the relationship was created
		records, err := result.Collect(c.ctx)
		if err != nil {
			return nil, err
		}

		if len(records) > 0 {
			return records[0].Values[0], nil
		}

		return nil, database.NewNotRelationshipCreatedError(userPair.FollowerID, userPair.FolloweeID)
	})

	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to create follow relationship")
		return err
	}

	return nil
}

func (c *Neo4jDBClient) DeleteRelationship(userPair *model.UserPairRelationship) error {
	driver, err := c.getNeo4jDriver()
	if err != nil {
		return err
	}
	defer driver.Close(c.ctx)
	session := driver.NewSession(c.ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(c.ctx)

	_, err = session.ExecuteWrite(c.ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (follower:User {id: $followerId})-[r:FOLLOWS]->(followee:User {id: $followeeId})
			DELETE r
		`
		params := map[string]interface{}{
			"followerId": userPair.FollowerID,
			"followeeId": userPair.FolloweeID,
		}

		result, err := tx.Run(c.ctx, query, params)
		if err != nil {
			return nil, err
		}

		summary, err := result.Consume(c.ctx)
		if err != nil {
			return nil, err
		}

		if summary.Counters().RelationshipsDeleted() == 0 {
			return nil, database.NewNotRelationshipDeletedError(userPair.FollowerID, userPair.FolloweeID)
		}

		return nil, nil
	})

	if err != nil {
		log.Error().Stack().Err(err).Msg("Failed to delete follow relationship")
		return err
	}

	return nil
}

func (c *Neo4jDBClient) Clean() {
	driver, err := c.getNeo4jDriver()
	if err != nil {
		return
	}
	defer driver.Close(c.ctx)
	session := driver.NewSession(c.ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(c.ctx)

	_, err = session.ExecuteWrite(c.ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryRelationships := `
            MATCH ()-[r:FOLLOWS]->()
            DELETE r
            RETURN count(r) as deletedRelationships
        `

		_, err := tx.Run(c.ctx, queryRelationships, nil)
		if err != nil {
			log.Error().Stack().Err(err).Msg("Error deleting relationships")
			return nil, err
		}

		queryNodes := `
    		MATCH (u:User)
    		DELETE u
    		RETURN count(u) as deletedNodes
		`

		_, err = tx.Run(c.ctx, queryNodes, nil)
		if err != nil {
			log.Error().Stack().Err(err).Msg("Error deleting nodes")
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		log.Error().Stack().Err(err).Msg("Unexpected error cleaning the database")
	}
}

func (c *Neo4jDBClient) GetUserFollowers(username string, lastFollowerId string, limit int) ([]string, string, error) {
	driver, err := c.getNeo4jDriver()
	if err != nil {
		return []string{}, "", err
	}
	defer driver.Close(c.ctx)

	query := `
		MATCH (followee:User {id: $followeeId})<-[:FOLLOWS]-(follower:User)
		WHERE ($lastFollowerId IS NULL OR follower.id > $lastFollowerId)
		RETURN follower.id AS followerId
		ORDER BY follower.id
		LIMIT $limit
	`
	params := map[string]interface{}{
		"followeeId":     username,
		"lastFollowerId": lastFollowerId,
		"limit":          limit,
	}

	session := driver.NewSession(c.ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(c.ctx)

	result, err := session.Run(c.ctx, query, params)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Get Followers Query exectution failed")
		return []string{}, "", err
	}

	// Processar os resultados
	var followers []string
	var lastId string

	for result.Next(c.ctx) {
		record := result.Record()
		followerId, ok := record.Get("followerId")
		if !ok {
			continue
		}

		followerIdStr, ok := followerId.(string)
		if !ok {
			continue
		}

		followers = append(followers, followerIdStr)
		lastId = followerIdStr
	}

	if err = result.Err(); err != nil {
		log.Error().Stack().Err(err).Msg("Error while processing results")
		return []string{}, "", err
	}

	return followers, lastId, nil
}

func (c *Neo4jDBClient) getNeo4jDriver() (neo4j.DriverWithContext, error) {
	driver, err := neo4j.NewDriverWithContext(
		c.dbUri,
		neo4j.BasicAuth(c.dbUser, c.dbPassword, ""))
	if err != nil {
		log.Error().Stack().Err(err).Msg("Neo4j driver was not created")
	}

	return driver, err
}
