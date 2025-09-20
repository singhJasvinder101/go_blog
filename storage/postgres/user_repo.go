package postgres

import (
	"context"
	"fmt"

	"github.com/singhJasvinder101/go_blog/internal/types"
)

type UserRepo struct {
	DB *Postgres
}

func NewUserRepo(db *Postgres) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}


func (u *UserRepo) CreateUser(ctx context.Context, user *types.User) (int, error) {
	query := `
	insert into users (name, email, password_hash) 
	values ($1, $2, $3) 
	returning id
	`

	row := u.DB.DB.QueryRow(ctx, query, user.Name, user.Email, user.PasswordHash)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert user error (user_repo) %w ", err)
	}
	return id, nil
}



func (u *UserRepo) GetUserByID(ctx context.Context, id int) (*types.User, error) {
	query := `
	select * from users where id = $1
	`
	row := u.DB.DB.QueryRow(ctx, query, id)

	var user types.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error when getting user by ID (user_repo) %w", err)

	}
	return &user, nil
}


