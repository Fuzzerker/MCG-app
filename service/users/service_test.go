package users

import (
	"context"
	"errors"
	"fmt"
	"mcg-app-backend/service/customerrors"
	"mcg-app-backend/service/models"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MockUsersRepo struct {
	mock.Mock
}

func (m *MockUsersRepo) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUsersRepo) InsertUser(ctx context.Context, user models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockTracer struct {
	mock.Mock
}

func (m *MockTracer) NewSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer("").Start(context.Background(), "")
}

func (m *MockTracer) SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {

}

func (m *MockTracer) RecordError(ctx context.Context, err error) error {
	return err
}

func getMocksAndService() (*MockUsersRepo, Service) {
	mockRepo := new(MockUsersRepo)
	mockTracer := new(MockTracer)
	service := NewService(mockRepo, mockTracer)
	return mockRepo, service
}
func TestCreateUser(t *testing.T) {
	username := "testuser"
	password := "password123"

	t.Run("CreateUser_Success", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetUserByUsername", mock.Anything, username).Return(models.User{}, nil)
		mockRepo.On("InsertUser", mock.Anything, mock.AnythingOfType("models.User")).Return(nil)

		err := service.CreateUser(context.Background(), username, password)
		assert.Nil(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("CreateUser_UsernameAlreadyExists", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetUserByUsername", mock.Anything, username).Return(models.User{Username: username}, nil)

		err := service.CreateUser(context.Background(), username, password)
		assert.NotNil(t, err)
		var alreadyExists customerrors.AlreadyExistsError
		assert.True(t, errors.As(err, &alreadyExists))

		mockRepo.AssertExpectations(t)
	})

	// Test repo error scenarios
	t.Run("CreateUser_RepoError_GetUserByUsername", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetUserByUsername", mock.Anything, username).Return(models.User{}, fmt.Errorf("db error"))

		err := service.CreateUser(context.Background(), username, password)
		assert.NotNil(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CreateUser_RepoError_InsertUser", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetUserByUsername", mock.Anything, username).Return(models.User{}, nil)
		mockRepo.On("InsertUser", mock.Anything, mock.AnythingOfType("models.User")).Return(fmt.Errorf("db error"))

		err := service.CreateUser(context.Background(), username, password)
		assert.NotNil(t, err)
		assert.Equal(t, "error inserting user db error", err.Error())

		mockRepo.AssertExpectations(t)
	})

}

func TestGetPasswordByUserName(t *testing.T) {

	username := "testuser"
	password := "hashedpassword123"
	t.Run("GetPasswordByUserName_Success", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetUserByUsername", mock.Anything, username).Return(models.User{Username: username, Password: password}, nil)

		retrievedPassword, err := service.GetPasswordByUserName(context.Background(), username)
		assert.Nil(t, err)
		assert.Equal(t, password, retrievedPassword)

		mockRepo.AssertExpectations(t)
	})
	t.Run("GetPasswordByUserName_UserNotFound", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetUserByUsername", mock.Anything, username).Return(models.User{}, nil)

		retrievedPassword, err := service.GetPasswordByUserName(context.Background(), username)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid username", err.Error())
		assert.Empty(t, retrievedPassword)

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetPasswordByUserName_ErrorFetchingUser", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetUserByUsername", mock.Anything, username).Return(models.User{}, fmt.Errorf("db error"))

		retrievedPassword, err := service.GetPasswordByUserName(context.Background(), username)
		assert.NotNil(t, err)
		assert.Equal(t, "error getting user by username db error", err.Error())
		assert.Empty(t, retrievedPassword)

		mockRepo.AssertExpectations(t)
	})
}
