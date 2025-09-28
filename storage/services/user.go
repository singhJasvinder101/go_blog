// this file contains business logic (cache before db call)
// creates abstraction over repository layer (DB implementation)

package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/singhJasvinder101/go_blog/internal/types"

	go_redis "github.com/redis/go-redis/v9"
)

//tests must not depend on concrete objects instances instead 
//manually write all the interfaces of that object and fake implement each
type UserService struct {
	UserRepo    UserRepository
	PostRepo    PostRepository
	RedisClient Cache
}

// DO NOT EXPOSE REDIS ERRORS INSTEAD EXPOSE ONLY UNEXPECTED ERRORS

func NewUserService(UserRepo UserRepository, postRepo PostRepository, redisClient Cache	) *UserService {
	return &UserService{
		UserRepo:    UserRepo,
		PostRepo:    postRepo,
		RedisClient: redisClient,
	}
}

func (s *UserService) Create(ctx context.Context, name, email, password string) (*types.User, error) {
	user := &types.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(password),
	}

	createdUser, err := s.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error while creating user (user_service) %w", err)
	}

	user_created_cache_key := fmt.Sprintf("user:%d", createdUser.ID)
	data, _ := json.Marshal(createdUser)

	s.RedisClient.Set(ctx, user_created_cache_key, data, 0)

	return createdUser, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*types.User, error) {
	user_cache_key := fmt.Sprintf("user:%d", id)

	data, err := s.RedisClient.Get(ctx, user_cache_key)
	// (err == go_redis.Nil)  => cache miss


	if err != nil && err != go_redis.Nil {
        return nil, err // exposing only unexpected error
    }

	// cache hit
	if err == nil { 
        var user types.User
        if unmarshalErr := json.Unmarshal([]byte(data), &user); unmarshalErr == nil {
            fmt.Println("cache hit")
            return &user, nil
        }
    }

	user, err := s.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		return user, nil
	}

	marshaledData, _ := json.Marshal(user)
	s.RedisClient.Set(ctx, user_cache_key, marshaledData, 5*time.Minute)
	
	return user, nil
}

func (s *UserService) GetByUID(ctx context.Context, userId int) ([]types.Post, error) {
	user_cache_key := fmt.Sprintf("user_posts:%d", userId)

	data, err := s.RedisClient.Get(ctx, user_cache_key)

	if err != nil && err != go_redis.Nil {
        return nil, err 
    }

	// cache hit 
    if err == nil { 
        var posts []types.Post
        if unmarshalErr := json.Unmarshal([]byte(data), &posts); unmarshalErr == nil {
            fmt.Println("cache hit")
            return posts, nil
        }
    }

	// cache miss
	posts, err := s.UserRepo.GetUserPosts(ctx, userId)
	if err != nil {
		return nil, err
	}

	marshaledData, _ := json.Marshal(posts)
	s.RedisClient.Set(ctx, user_cache_key, marshaledData, 5*time.Minute)

	return posts, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*types.User, error) {
	user_cache_key := fmt.Sprintf("user_get_by_email:%s", email)

	data, err := s.RedisClient.Get(ctx, user_cache_key)
	
	if err != nil && err != go_redis.Nil {
		return nil, err 
	}

	// cache hit
	if err == nil { 
		var user types.User
		if unmarshalErr := json.Unmarshal([]byte(data), &user); unmarshalErr == nil {
			fmt.Println("cache hit")
			return &user, nil
		}
	}
	
	user, err := s.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	
	marshaledData, _ := json.Marshal(user)
	s.RedisClient.Set(ctx, user_cache_key, marshaledData, 5*time.Minute)

	return user, nil
}

func (s *UserService) CreateComment(ctx context.Context, content string, userId, postId int) (types.Comment, error) {
	if content == "" {
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
		UserID:  userId,
		PostID:  postId,
	}

	fmt.Print(comment)
	return s.UserRepo.AddComment(ctx, comment)
}
