package services

import (
	"context"

	"github.com/singhJasvinder101/go_blog/internal/types"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *types.Post) (int, error)
	GetAllPosts(ctx context.Context) ([]types.Post, error)
	GetPostByID(ctx context.Context, postId int) (types.Post, error)
}



