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

func (s *UserService) CreateUser(name, email, password string) (int, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("failed to hash password")
	}
	
	user := &types.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashed),
	}
	return s.UserRepo.CreateUser(context.Background(), user)
}

func (s *UserService) GetUserByID(id int) (*types.User, error) {
	return s.UserRepo.GetUserByID(context.Background(), id)
}	


