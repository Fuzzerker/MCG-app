package attatchments

import (
	"context"
	"fmt"
	"mcg-app-backend/service/customerrors"
	"mcg-app-backend/service/models"

	"go.opentelemetry.io/otel/attribute"
)

type AttachmentService struct {
	repo       AttachmentRepo
	patientSvc PatientService
	tracer     Tracer
}

func NewAttachmentService(repo AttachmentRepo, patientSvc PatientService, tracer Tracer) AttachmentService {
	return AttachmentService{
		repo:       repo,
		patientSvc: patientSvc,
		tracer:     tracer,
	}
}

func (s AttachmentService) AddAttatchmentToPatient(ctx context.Context, patientId int, name string, description string, typ string, data []byte) (models.Attatchment, error) {
	ctx, span := s.tracer.NewSpan(ctx, "AddAttachmentToPatient")
	defer span.End()
	s.tracer.SetAttributes(ctx,
		attribute.Int("patientId", patientId),
		attribute.String("name", name),
		attribute.String("description", description),
		attribute.String("type", typ))

	if len(data) == 0 {
		return models.Attatchment{}, s.tracer.RecordError(ctx, customerrors.NewInvalidInputError("data was empty"))
	}

	err := s.patientSvc.ValidatePatientId(ctx, patientId)
	if err != nil {
		return models.Attatchment{}, err
	}

	attachment := models.Attatchment{
		PatientId:   patientId,
		Name:        name,
		Description: description,
		Type:        typ,
		Data:        data,
	}

	id, err := s.repo.InsertAttatchment(ctx, attachment)
	if err != nil {
		return models.Attatchment{}, s.tracer.RecordError(ctx, fmt.Errorf("error inserting attachment %w", err))
	}

	attachment.Id = id
	return attachment, nil
}

func (s AttachmentService) DeleteAttatchment(ctx context.Context, attachmentId int) error {
	ctx, span := s.tracer.NewSpan(ctx, "DeleteAttachment")
	defer span.End()
	s.tracer.SetAttributes(ctx, attribute.Int("attachmentId", attachmentId))

	err := s.repo.DeleteAttatchment(ctx, attachmentId)
	if err != nil {
		return s.tracer.RecordError(ctx, fmt.Errorf("error deleting attachment %w", err))
	}

	return nil
}

func (s AttachmentService) DeletePatientAttachments(ctx context.Context, patientId int) error {
	ctx, span := s.tracer.NewSpan(ctx, "DeletePatientAttachments")
	defer span.End()
	s.tracer.SetAttributes(ctx, attribute.Int("patientId", patientId))

	err := s.repo.DeleteAttatchmentsByPatientId(ctx, patientId)
	if err != nil {
		return s.tracer.RecordError(ctx, fmt.Errorf("error deleting attachments for patient %w", err))
	}

	return nil
}
