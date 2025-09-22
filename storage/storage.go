package storage

import (
	"github.com/singhJasvinder101/go_blog/internal/config"
	"github.com/singhJasvinder101/go_blog/internal/types"
	"github.com/singhJasvinder101/go_blog/storage/postgres"
)

type Storage interface {
	NewPostgres(cfg *config.Config) (*postgres.Postgres, error)
	InitSchema() error
	Close() error
}

type UserService interface {
	CreateUser(name, email, password string) (int, error)
	GetUserByID(id string) (*types.User, error)
	GetUserPosts(userId int) ([]types.Post, error)
	CreateUserComment(content string, userId, postId int) (types.Comment, error)
}

type PostService interface {
	CreatePost(title, description string, userId int) (int, error)
	GetPostByID(id int) (*types.Post, error)
}



