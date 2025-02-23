package users

import (
	"context"
	"mcg-app-backend/service/models"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UsersRepo interface {
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	InsertUser(ctx context.Context, user models.User) error
}

type Tracer interface {
	NewSpan(ctx context.Context, name string) (context.Context, trace.Span)
	SetAttributes(ctx context.Context, attrs ...attribute.KeyValue)
	RecordError(ctx context.Context, err error) error
}
