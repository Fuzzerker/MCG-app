package attatchments

import (
	"context"
	"mcg-app-backend/service/models"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AttachmentRepo interface {
	InsertAttatchment(ctx context.Context, attachment models.Attatchment) (int, error)
	DeleteAttatchment(ctx context.Context, attachmentId int) error
	DeleteAttatchmentsByPatientId(ctx context.Context, patientId int) error
}

type PatientService interface {
	ValidatePatientId(ctx context.Context, patientId int) error
}

type Tracer interface {
	NewSpan(ctx context.Context, name string) (context.Context, trace.Span)
	SetAttributes(ctx context.Context, attrs ...attribute.KeyValue)
	RecordError(ctx context.Context, err error) error
}
