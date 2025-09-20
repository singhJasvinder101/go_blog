package postgres

import (
	"context"
	"fmt"

	"github.com/singhJasvinder101/go_blog/internal/types"
)

type PostRepo struct {
	DB *Postgres
}

func NewPostRepo(db *Postgres) *PostRepo {
	return &PostRepo{
		DB: db,
	}
}

func (r *PostRepo) CreatePost(ctx context.Context, post *types.Post) (int, error) {
	query := `
	insert into posts (title, description, user_id)
	values ($1, $2, $3)
	returning id
	`

	row := r.DB.Pool.QueryRow(ctx, query, post.Title, post.Description, post.UserID)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert post error (post_repo) %w ", err)
	}
	return id, nil
}
