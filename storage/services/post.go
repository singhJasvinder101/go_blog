// this file contains business logic (cache before db call)
// creates abstraction over repository layer (DB implementation)

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/singhJasvinder101/go_blog/internal/types"

	go_redis "github.com/redis/go-redis/v9"
)

type PostService struct {
	PostRepo   PostRepository
	RedisClient Cache
}

func NewPostService(postRepo PostRepository, redisClient Cache) *PostService{
	return &PostService{
		PostRepo:   postRepo,
		RedisClient: redisClient,
	}
}

func (s *PostService) Create(ctx context.Context, title, description string, userId int) (int, error) {
	post := types.Post{
		Title:       title,
		Description: description,
		UserID:      userId,
	}

	createdPostID, err := s.PostRepo.CreatePost(ctx, &post)
	if err != nil {
		return 0, err
	}

	post_cache_key := fmt.Sprintf("post:%d", createdPostID)
	marshaledData, _ := json.Marshal(createdPostID)
	s.RedisClient.Set(ctx, post_cache_key, marshaledData, 5*time.Minute)

	return createdPostID, nil
}

func (s *PostService) GetByID(ctx context.Context, postId int) (types.Post, error){
	post_cache_key := fmt.Sprintf("post:%d", postId)
	
	data, err := s.RedisClient.Get(ctx, post_cache_key)
	// (err == go_redis.Nil)  => cache miss
	
	if err != nil && err != go_redis.Nil {
		return types.Post{}, err
	}
	
	// cache hit
	if err == nil {
		var post types.Post
		err = json.Unmarshal([]byte(data), &post)
		if err != nil {
			return types.Post{}, err
		}
		return post, nil
	}

	post, err := s.PostRepo.GetPostByID(ctx, postId)
	if err != nil {
		return types.Post{}, err
	}
	
	marshaledData, _ := json.Marshal(post)
	s.RedisClient.Set(ctx, post_cache_key, marshaledData, 5*time.Minute)

	return post, nil
}

func (s *PostService) GetAll(ctx context.Context) ([]types.Post, error) {
	posts_cache_key := "posts:all"

	data, err := s.RedisClient.Get(ctx, posts_cache_key)
	// (err == go_redis.Nil)  => cache miss

	if err != nil && err != go_redis.Nil {
		return nil, err
	}

	// cache hit
	if err == nil {
		var posts []types.Post
		err = json.Unmarshal([]byte(data), &posts)
		if err != nil {
			return nil, err
		}
		return posts, nil
	}

	posts, err := s.PostRepo.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	marshaledData, _ := json.Marshal(posts)
	s.RedisClient.Set(ctx, posts_cache_key, marshaledData, 5*time.Minute)

	return posts, nil
}

