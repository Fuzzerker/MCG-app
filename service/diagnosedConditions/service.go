package diagnosedconditions

import (
	"context"
	"fmt"
	"mcg-app-backend/service/models"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

type DiagnosedConditionService struct {
	repo       DiagnosedConditionRepo
	patientSvc PatientService
	tracer     Tracer
}

func NewDiagnosedConditionService(repo DiagnosedConditionRepo, patientSvc PatientService, tracer Tracer) DiagnosedConditionService {
	return DiagnosedConditionService{
		repo:       repo,
		patientSvc: patientSvc,
		tracer:     tracer,
	}
}

func (s DiagnosedConditionService) AddDiagnosedConditionToPatient(ctx context.Context, patientId int, name string, code string, description string, date time.Time) (models.DiagnosedCondition, error) {
	ctx, span := s.tracer.NewSpan(ctx, "AddDiagnosedConditionToPatient")
	defer span.End()
	s.tracer.SetAttributes(ctx,
		attribute.Int("patientId", patientId),
		attribute.String("name", name),
		attribute.String("code", code),
		attribute.String("description", description),
		attribute.String("date", fmt.Sprintf("%v", date)))

	err := s.patientSvc.ValidatePatientId(ctx, patientId)
	if err != nil {
		return models.DiagnosedCondition{}, err
	}

	condition := models.DiagnosedCondition{
		PatientId:   patientId,
		Name:        name,
		Code:        code,
		Description: description,
		Date:        date,
	}

	id, err := s.repo.InsertDiagnosedCondition(ctx, condition)
	if err != nil {
		return models.DiagnosedCondition{}, s.tracer.RecordError(ctx, fmt.Errorf("error inserting diagnosed condition %w", err))
	}

	condition.Id = id
	return condition, nil
}

func (s DiagnosedConditionService) DeleteDiagnosedCondition(ctx context.Context, conditionId int) error {
	ctx, span := s.tracer.NewSpan(ctx, "DeleteDiagnosedCondition")
	defer span.End()
	s.tracer.SetAttributes(ctx, attribute.Int("conditionId", conditionId))

	err := s.repo.DeleteDiagnosedCondition(ctx, conditionId)
	if err != nil {
		return s.tracer.RecordError(ctx, fmt.Errorf("error deleting diagnosed condition %w", err))
	}

	return nil
}

func (s DiagnosedConditionService) DeletePatientDiagnosedConditions(ctx context.Context, patientId int) error {
	ctx, span := s.tracer.NewSpan(ctx, "DeletePatientDiagnosedConditions")
	defer span.End()
	s.tracer.SetAttributes(ctx, attribute.Int("patientId", patientId))

	err := s.repo.DeleteDiagnosedConditionsByPatientId(ctx, patientId)
	if err != nil {
		return s.tracer.RecordError(ctx, fmt.Errorf("error deleting diagnosed conditions for patient %w", err))
	}

	return nil
}
