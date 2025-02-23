package diagnosedconditions

import (
	"context"
	"mcg-app-backend/service/models"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DiagnosedConditionRepo interface {
	InsertDiagnosedCondition(ctx context.Context, condition models.DiagnosedCondition) (int, error)
	DeleteDiagnosedCondition(ctx context.Context, conditionId int) error
	DeleteDiagnosedConditionsByPatientId(ctx context.Context, patientId int) error
}

type PatientService interface {
	ValidatePatientId(ctx context.Context, patientId int) error
}

type Tracer interface {
	NewSpan(ctx context.Context, name string) (context.Context, trace.Span)
	SetAttributes(ctx context.Context, attrs ...attribute.KeyValue)
	RecordError(ctx context.Context, err error) error
}
