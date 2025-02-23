package auth

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UsersService interface {
	GetPasswordByUserName(ctx context.Context, username string) (string, error)
}

type Tracer interface {
	NewSpan(ctx context.Context, name string) (context.Context, trace.Span)
	SetAttributes(ctx context.Context, attrs ...attribute.KeyValue)
	RecordError(ctx context.Context, err error) error
}
