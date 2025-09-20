package storage

import "github.com/singhJasvinder101/go_blog/internal/types"

type Storage interface {
	InitSchema() error
}

type UserService interface {
	CreateUser(name, email, password string) (int, error)
	GetUserByID(id string) (*types.User, error)
}


