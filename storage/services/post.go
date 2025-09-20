package services

import (
	"context"

	"github.com/singhJasvinder101/go_blog/internal/types"
	"github.com/singhJasvinder101/go_blog/storage/postgres"
)

type PostService struct {
	PostRepo *postgres.PostRepo
}

func NewPostService(postRepo *postgres.PostRepo) *PostService{
	return &PostService{
		PostRepo: postRepo,
	}
}

func (s *PostService) Create(ctx context.Context, title, description string, userId int) (int, error) {
	post := types.Post{
		Title:       title,
		Description: description,
		UserID:      userId,
	}
	return s.PostRepo.CreatePost(ctx, &post)
}

func (s *PostService) GetByID(ctx context.Context, postId int) (types.Post, error){
	return s.PostRepo.GetPostByID(ctx, postId)
}

func (s *PostService) GetAll(ctx context.Context) ([]types.Post, error) {
	return s.PostRepo.GetAllPosts(ctx)
}

