package services

import (
	"context"
	"time"

	"github.com/singhJasvinder101/go_blog/internal/types"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *types.User) (*types.User, error)
	GetUserByID(ctx context.Context, id int) (*types.User, error)
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	GetUserPosts(ctx context.Context, id int) ([]types.Post, error)
	AddComment(ctx context.Context, c types.Comment) (types.Comment, error)
}


type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val interface{}, exp time.Duration) error
}
