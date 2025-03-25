package provider

import (
	"context"
	dbInfra "followservice/infrastructure/database"
	"followservice/infrastructure/kafka"
	"followservice/internal/api"
	"followservice/internal/bus"
	database "followservice/internal/db"
	"followservice/internal/features/follow_user"
	"followservice/internal/features/get_user_followers"
	"followservice/internal/features/unfollow_user"
)

type Provider struct {
	env string
}

func NewProvider(env string) *Provider {
	return &Provider{
		env: env,
	}
}

func (p *Provider) ProvideEventBus() (*bus.EventBus, error) {
	kafkaProducer, err := kafka.NewKafkaProducer(p.kafkaBrokers())
	if err != nil {
		return nil, err
	}

	return bus.NewEventBus(kafkaProducer), nil
}

func (p *Provider) ProvideSubscriptions() *[]bus.EventSubscription {
	return &[]bus.EventSubscription{
		{},
	}
}

func (p *Provider) ProvideKafkaConsumer(eventBus *bus.EventBus) (*kafka.KafkaConsumer, error) {
	brokers := p.kafkaBrokers()

	return kafka.NewKafkaConsumer(brokers, eventBus)
}

func (p *Provider) ProvideApiEndpoint(database *database.Database, cache *database.Cache, bus *bus.EventBus) *api.Api {
	return api.NewApiEndpoint(p.env, p.ProvideApiControllers(database, cache, bus))
}

func (p *Provider) ProvideApiControllers(database *database.Database, cache *database.Cache, bus *bus.EventBus) []api.Controller {
	return []api.Controller{
		follow_user.NewFollowUserController(follow_user.NewFollowUserService(follow_user.NewFollowUserRepository(database), bus)),
		unfollow_user.NewUnfollowUserController(unfollow_user.NewUnfollowUserService(unfollow_user.NewUnfollowUserRepository(database), bus)),
		get_user_followers.NewGetUserFollowersController(get_user_followers.NewGetUserFollowersService(get_user_followers.NewGetUserFollowersRepository(database, cache))),
	}
}

func (p *Provider) ProvideDb(ctx context.Context) *database.Database {
	return database.NewDatabase(dbInfra.NewNeo4jClient("bolt://localhost:7687", "neo4j", "contrasinal", ctx))
}

func (p *Provider) ProvideCache(ctx context.Context) *database.Cache {
	return database.NewCache(dbInfra.NewRedisClient("localhost:6379", "", ctx))
}

func (p *Provider) kafkaBrokers() []string {
	if p.env == "development" || p.env == "test" {
		return []string{
			"localhost:9093",
		}
	} else {
		return []string{
			"172.31.36.175:9092",
			"172.31.45.255:9092",
		}
	}
}
