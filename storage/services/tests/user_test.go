package services_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/singhJasvinder101/go_blog/internal/types"
	"github.com/singhJasvinder101/go_blog/storage/services"
)

// ---- MOCKS ----

type (
	MockUserRepo    struct{ mock.Mock }
	MockPostRepo    struct{ mock.Mock }
	MockRedisClient struct{ mock.Mock }
)

func (m *MockUserRepo) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, id int) (*types.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*types.User), args.Error(1)
}

func (m *MockUserRepo) GetUserPosts(ctx context.Context, id int) ([]types.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]types.Post), args.Error(1)
}

func (m *MockUserRepo) AddComment(ctx context.Context, c types.Comment) (types.Comment, error) {
	args := m.Called(ctx, c)
	return args.Get(0).(types.Comment), args.Error(1)
}

func (m *MockPostRepo) GetPostByID(ctx context.Context, id int) (types.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(types.Post), args.Error(1)
}

func (m *MockPostRepo) CreatePost(ctx context.Context, post *types.Post) (int, error) {
	args := m.Called(ctx, post)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockPostRepo) GetAllPosts(ctx context.Context) ([]types.Post, error) {
	args := m.Called(ctx)
	return args.Get(0).([]types.Post), args.Error(1)
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, val interface{}, exp time.Duration) error {
	args := m.Called(ctx, key, val, exp)
	return args.Error(0)
}

// ---- TESTS ----

func TestUserService_Create(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := new(MockUserRepo)
	mockPostRepo := new(MockPostRepo)
	mockRedis := new(MockRedisClient)

	service := services.NewUserService(mockUserRepo, mockPostRepo, mockRedis)

	expectedUser := &types.User{
		ID:           1,
		Name:         "John",
		Email:        "j@doe.com",
	}

	// userRepo.CreateUser returns expected user (Fake DB behavior)
	mockUserRepo.On("CreateUser", ctx, mock.AnythingOfType("*types.User")).Return(expectedUser, nil)
	mockRedis.On("Set", ctx, "user:1", mock.AnythingOfType("[]uint8"), 0*time.Second).Return(nil)

	user, err := service.Create(ctx, "John", "j@doe.com", "pass")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserService_GetByID_CacheHit(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := new(MockUserRepo)
	mockPostRepo := new(MockPostRepo)
	mockRedis := new(MockRedisClient)

	service := services.NewUserService(mockUserRepo, mockPostRepo, mockRedis)

	expectedUser := &types.User{
		ID:           1,
		Name:         "John",
		Email:        "j@doe.com",
	}

	data, _ := json.Marshal(expectedUser)
	mockRedis.On("Get", ctx, "user:1").Return(string(data), nil)

	user, err := service.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	// db should not be called
	mockUserRepo.AssertNotCalled(t, "GetUserByID", ctx, 1)
}
