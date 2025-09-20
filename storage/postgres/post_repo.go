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

func (r *PostRepo) GetPostByID(ctx context.Context, postId int) (types.Post, error){
	query := `
	select * from posts where id = $1
	`

	row := r.DB.Pool.QueryRow(ctx, query, postId)

	var post types.Post
	err := row.Scan(&post.ID, &post.Title, &post.Description, &post.UserID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return types.Post{}, fmt.Errorf("error when getting post by ID (post_repo) %w", err)
	}

	return post, nil
}

func (r *PostRepo) GetAllPosts(ctx context.Context) ([]types.Post, error){
	query := `
	select * from posts
	`

	rows, err := r.DB.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all posts error (post_repo) %w ", err)
	}
	defer rows.Close()
	
	var posts []types.Post
	for rows.Next() {
		var post types.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.UserID, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan post error (post_repo) %w ", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}
