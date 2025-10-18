//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/singhJasvinder101/go_blog/internal/config"
	posthandlers "github.com/singhJasvinder101/go_blog/internal/http/handlers/posts"
	userhandlers "github.com/singhJasvinder101/go_blog/internal/http/handlers/users"
	"github.com/singhJasvinder101/go_blog/storage/postgres"
	"github.com/singhJasvinder101/go_blog/storage/redis"
	"github.com/singhJasvinder101/go_blog/storage/services"
)

func ProvideRedisClient(cfg *config.Config) (*redis.RedisClient, error) {
	return redis.NewRedisClient(cfg, 0)
}

//interface bindings (concrete types to interfaces)
func ProvideUserRepository(repo *postgres.UserRepo) services.UserRepository {
    return repo
}
func ProvidePostRepository(repo *postgres.PostRepo) services.PostRepository {
    return repo
}
func ProvideCache(client *redis.RedisClient) services.Cache {
    return client
}


// all providers ka set (constructors/New functions)
var ProviderSet = wire.NewSet(
	//database providers
	postgres.NewPostgres,
	ProvideRedisClient,

	//repos
	postgres.NewUserRepo,
	postgres.NewPostRepo,

	//interface bindings
	ProvideUserRepository,
	ProvidePostRepository,
	ProvideCache,

	//services
	services.NewUserService,
	services.NewPostService,

	//handlers
	userhandlers.NewUserHandler,
	posthandlers.NewPostHandler,
)


type Application struct {
	UserHandler *userhandlers.UserHandler
	PostHandler *posthandlers.PostHandler
	DB          *postgres.Postgres
	Redis       *redis.RedisClient
}

func InitializeApplication(cfg *config.Config) (*Application, error) {
	wire.Build(
		ProviderSet,
		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil
}
