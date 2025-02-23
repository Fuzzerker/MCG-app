package auth

import (
	"context"
	"fmt"
	"mcg-app-backend/service/customerrors"
	"mcg-app-backend/service/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	usersService        UsersService
	tracer              Tracer
	tokenExpirationTime time.Duration
	issuer              string
	tokenSecretKey      []byte
}

func NewService(usersService UsersService, tracer Tracer, tokenExpirationTime time.Duration,
	issuer string, tokenSecretKey string) Service {
	return Service{
		usersService:        usersService,
		tracer:              tracer,
		tokenExpirationTime: tokenExpirationTime,
		issuer:              issuer,
		tokenSecretKey:      []byte(tokenSecretKey),
	}
}

func (s Service) Login(ctx context.Context, username string, password string) (string, error) {
	ctx, span := s.tracer.NewSpan(ctx, "Login")
	defer span.End()
	s.tracer.SetAttributes(ctx, attribute.String("username", username))
	storedPassword, err := s.usersService.GetPasswordByUserName(ctx, username)
	if err != nil {
		return "", s.tracer.RecordError(ctx, fmt.Errorf("error getting user from db %w", err))
	}

	compErr := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if compErr != nil {
		s.tracer.RecordError(ctx, fmt.Errorf("error in bcrypt compare %w", compErr))
		return "", s.tracer.RecordError(ctx, customerrors.NewInvalidInputError("password does not match"))
	}

	token, err := s.generateToken(username)
	if err != nil {
		return "", s.tracer.RecordError(ctx, fmt.Errorf("error generating token %w", err))
	}
	return token, nil
}

func (s Service) generateToken(username string) (string, error) {
	expiration := time.Now().Add(s.tokenExpirationTime)

	claims := models.UserClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.tokenSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (s Service) VerifyToken(ctx context.Context, tokenString string) error {
	ctx, span := s.tracer.NewSpan(ctx, "VerifyToken")
	defer span.End()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.tokenSecretKey, nil
	})
	if err != nil {
		return s.tracer.RecordError(ctx, fmt.Errorf("error parsing token %w", err))
	}
	s.tracer.SetAttributes(ctx, attribute.Bool("token.valid", token.Valid))
	if !token.Valid {
		return s.tracer.RecordError(ctx, customerrors.NewUnauthorizedError("token is invalid"))
	}

	return nil
}
