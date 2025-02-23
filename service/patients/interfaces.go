package patients

import (
	"context"
	"mcg-app-backend/service/models"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type PatientRepo interface {
	GetCountOfExternalIdentifier(ctx context.Context, externalId string) (int, error)
	GetCountOfPatientId(ctx context.Context, patientId int) (int, error)
	InsertPatient(ctx context.Context, patient models.Patient) (int, error)
	UpdatePatient(ctx context.Context, patient models.Patient) error
	DeletePatient(ctx context.Context, patientId int) error
	SearchPatients(ctx context.Context, search models.PatientSearch) ([]models.Patient, error)
}

type Tracer interface {
	NewSpan(ctx context.Context, name string) (context.Context, trace.Span)
	SetAttributes(ctx context.Context, attrs ...attribute.KeyValue)
	RecordError(ctx context.Context, err error) error
}
