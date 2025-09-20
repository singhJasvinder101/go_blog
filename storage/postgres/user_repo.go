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

func (r *UserRepo) CreateUser(ctx context.Context, user *types.User) (int, error) {
	query := `
	insert into users (name, email, password_hash) 
	values ($1, $2, $3) 
	returning id
	`

	row := r.DB.Pool.QueryRow(ctx, query, user.Name, user.Email, user.PasswordHash)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert user error (user_repo) %w ", err)
	}
	return id, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int) (*types.User, error) {
	query := `
	select * from users where id = $1
	`
	row := r.DB.Pool.QueryRow(ctx, query, id)

	var user types.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error when getting user by ID (user_repo) %w", err)
	}
	return &user, nil
}


func (r *UserRepo) GetUserPosts(ctx context.Context, userId int) ([]types.Post, error) {
	query := `
	select * from posts where user_id = $1
	`

	rows, err := r.DB.Pool.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("error when getting posts by user ID (user_repo) %w", err)
	}

	var posts []types.Post
	
	for rows.Next(){
		var post types.Post

		err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error when scanning post row (user_repo) %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *UserRepo) AddComment(ctx context.Context, comment types.Comment) (types.Comment, error){
	query := `
	insert into comments (content, user_id, post_id)
	values ($1, $2, $3)
	returning id, content, user_id, post_id, created_at
	`

	row := r.DB.Pool.QueryRow(ctx, query, comment.Content, comment.UserID, comment.PostID)

	var val types.Comment
	err := row.Scan(&val.ID, &val.Content, &val.UserID, &val.PostID, &val.CreatedAt)

	if err != nil {
		return types.Comment{}, fmt.Errorf("error when adding comment %w", err)
	}

	return val, nil
}

