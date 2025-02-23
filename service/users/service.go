package users

import (
	"context"
	"fmt"
	"mcg-app-backend/service/customerrors"
	"mcg-app-backend/service/models"

	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   UsersRepo
	tracer Tracer
}

func NewService(repo UsersRepo, tracer Tracer) Service {
	return Service{
		repo:   repo,
		tracer: tracer,
	}
}

func (s Service) GetPasswordByUserName(ctx context.Context, username string) (string, error) {
	ctx, span := s.tracer.NewSpan(ctx, "GetPasswordByUserName")
	defer span.End()
	s.tracer.SetAttributes(ctx, attribute.String("username", username))
	user, err := s.repo.GetUserByUsername(ctx, username)

	if err != nil {
		return "", s.tracer.RecordError(ctx, fmt.Errorf("error getting user by username %w", err))
	}
	if user.Password == "" {
		return "", s.tracer.RecordError(ctx, customerrors.NewInvalidInputError("invalid username"))
	}
	return user.Password, nil
}

func (s Service) CreateUser(ctx context.Context, username string, password string) error {
	ctx, span := s.tracer.NewSpan(ctx, "CreateUser")
	defer span.End()
	s.tracer.SetAttributes(ctx, attribute.String("username", username))
	err := s.validateUniqueUsername(ctx, username)
	if err != nil {
		return s.tracer.RecordError(ctx, fmt.Errorf("error validating username %w", err))
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return s.tracer.RecordError(ctx, fmt.Errorf("error hashing password %w", err))
	}

	user := models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	err = s.repo.InsertUser(ctx, user)
	if err != nil {
		return s.tracer.RecordError(ctx, fmt.Errorf("error inserting user %w", err))
	}
	return nil
}

func (s Service) validateUniqueUsername(ctx context.Context, username string) error {
	existingUser, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("error getting user by username %w", err)
	}
	if existingUser.Username == username {
		return customerrors.NewAlreadyExistsError("Username is already taken")
	}
	return nil
}
