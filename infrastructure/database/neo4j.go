package database

import (
	"context"

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
	driver, err := neo4j.NewDriverWithContext(
		c.dbUri,
		neo4j.BasicAuth(c.dbUser, c.dbPassword, ""))
	if err != nil {
		log.Error().Stack().Err(err).Msg("Neo4j driver was not created")
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
