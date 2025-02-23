package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service struct {
	tracer trace.Tracer
	logger *zap.Logger
}

func NewService(logger *zap.Logger) Service {
	return Service{
		tracer: otel.Tracer("mcg-app"),
		logger: logger,
	}
}

func (s Service) NewSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return s.tracer.Start(ctx, name)
}

func (s Service) RecordError(ctx context.Context, err error) error {
	zap.Error(err)
	trace.SpanFromContext(ctx).RecordError(err)
	return err
}

func (s Service) SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	trace.SpanFromContext(ctx).SetAttributes(attrs...)
}
