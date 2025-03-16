package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisCacheClient struct {
	ctx    context.Context
	client *redis.Client
}

func NewRedisClient(cacheUri, cachePassword string, ctx context.Context) *RedisCacheClient {
	redisConfig := &redis.Options{
		Addr:     cacheUri,
		Password: cachePassword,
		DB:       0, // Use default DB
		Protocol: 2, // Connection protocol
	}
	client := &RedisCacheClient{
		ctx:    ctx,
		client: redis.NewClient(redisConfig),
	}

	client.verifyConnection()

	return client
}

func (c *RedisCacheClient) verifyConnection() {
	err := c.client.Set(c.ctx, "foo", "bar", 10*time.Second).Err()
	if err != nil {
		log.Error().Stack().Err(err).Msg("Conection to Redis not stablished")
		panic(err)
	}

	_, err = c.client.Get(c.ctx, "foo").Result()
	if err != nil {
		log.Error().Stack().Err(err).Msg("Conection to Redis not stablished")
		panic(err)
	}
	log.Info().Msgf("Connection to Redis established.")
}
