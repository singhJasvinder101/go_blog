package services

import (
	"context"
	"errors"

	"github.com/singhJasvinder101/go_blog/internal/types"
	"github.com/singhJasvinder101/go_blog/storage/postgres"
	"golang.org/x/crypto/bcrypt"
)


type UserService struct {
	UserRepo *postgres.UserRepo
}

func NewUserService(UserRepo *postgres.UserRepo) *UserService{
	return &UserService{
		UserRepo: UserRepo,
	}
}

func (s *UserService) Create(ctx context.Context,name, email, password string) (int, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("failed to hash password")
	}
	
	user := &types.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashed),
	}
	return s.UserRepo.CreateUser(ctx, user)
}

func (s *UserService) GetByID(ctx context.Context, id int) (*types.User, error) {
	return s.UserRepo.GetUserByID(ctx, id)
}	

func (s *UserService) GetByUID(ctx context.Context, userId int) ([]types.Post, error) {
	return s.UserRepo.GetUserPosts(ctx, userId)
}


