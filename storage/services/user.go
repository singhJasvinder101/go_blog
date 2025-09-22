package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/singhJasvinder101/go_blog/internal/types"
	"github.com/singhJasvinder101/go_blog/storage/postgres"
)


type UserService struct {
	UserRepo *postgres.UserRepo
	PostRepo *postgres.PostRepo
}

func NewUserService(UserRepo *postgres.UserRepo, postRepo *postgres.PostRepo) *UserService{
	return &UserService{
		UserRepo: UserRepo,
		PostRepo: postRepo,
	}
}

func (s *UserService) Create(ctx context.Context,name, email, password string) (*types.User, error) {
	
	user := &types.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(password),
	}
	return s.UserRepo.CreateUser(ctx, user)
}

func (s *UserService) GetByID(ctx context.Context, id int) (*types.User, error) {
	return s.UserRepo.GetUserByID(ctx, id)
}	

func (s *UserService) GetByUID(ctx context.Context, userId int) ([]types.Post, error) {
	return s.UserRepo.GetUserPosts(ctx, userId)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*types.User, error){
	return s.UserRepo.GetUserByEmail(ctx, email)	
}

func (s *UserService) CreateComment(ctx context.Context, content string, userId, postId int) (types.Comment, error){
	if content == ""{
		return types.Comment{}, errors.New("please add valid comment")
	}
	
	// GetByUID is method of this serivce module
	_, err := s.GetByUID(ctx, userId)
	if err != nil {
		return types.Comment{}, fmt.Errorf("invalid user id %w", err)
	}

	_, err = s.PostRepo.GetPostByID(ctx, postId)
	if err != nil {
		return types.Comment{}, fmt.Errorf("invalid post id %w", err)
	}

	comment := types.Comment{
		Content: content,
		UserID: userId,
		PostID: postId,
	}

	fmt.Print(comment)
	return s.UserRepo.AddComment(ctx, comment)
}


