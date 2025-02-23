// patient_service.go
package patients

import (
	"context"
	"fmt"
	"mcg-app-backend/service/customerrors"
	"mcg-app-backend/service/models"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

type PatientService struct {
	repo   PatientRepo
	tracer Tracer
}

func NewPatientService(repo PatientRepo, tracer Tracer) PatientService {
	return PatientService{
		repo:   repo,
		tracer: tracer,
	}
}

func (s PatientService) CreatePatient(ctx context.Context, name string, address string, phoneNumber string, dateOfBirth time.Time, externalIdentifier string) (models.Patient, error) {
	ctx, span := s.tracer.NewSpan(ctx, "CreatePatient")
	defer span.End()
	s.tracer.SetAttributes(ctx,
		attribute.String("address", address),
		attribute.String("name", name),
		attribute.String("dateOfBirth", fmt.Sprintf("%v", dateOfBirth)),
		attribute.String("phoneNumber", phoneNumber))

	count, err := s.repo.GetCountOfExternalIdentifier(ctx, externalIdentifier)
	if err != nil {
		return models.Patient{}, s.tracer.RecordError(ctx, fmt.Errorf("error getting count of externalIdentifier %w", err))
	}
	if count > 0 {
		return models.Patient{}, s.tracer.RecordError(ctx, customerrors.NewAlreadyExistsError("patient with matching externalIdentifier already exists"))
	}

	patient := models.Patient{
		Name:               name,
		PhoneNumber:        phoneNumber,
		Address:            address,
		ExternalIdentifier: externalIdentifier,
		DateOfBirth:        dateOfBirth,
	}

	id, err := s.repo.InsertPatient(ctx, patient)
	if err != nil {
		return models.Patient{}, s.tracer.RecordError(ctx, fmt.Errorf("error inserting patient %w", err))
	}
	patient.Id = id
	return patient, nil
}

func (s PatientService) UpdatePatient(ctx context.Context, id int, name string, address string, phoneNumber string, dateOfBirth time.Time, externalIdentifier string) (models.Patient, error) {
	ctx, span := s.tracer.NewSpan(ctx, "UpdatePatient")
	defer span.End()
	s.tracer.SetAttributes(ctx,
		attribute.Int("id", id),
		attribute.String("address", address),
		attribute.String("name", name),
		attribute.String("dateOfBirth", fmt.Sprintf("%v", dateOfBirth)),
		attribute.String("phoneNumber", phoneNumber),
	)

	err := s.ValidatePatientId(ctx, id)
	if err != nil {
		return models.Patient{}, err
	}

	patient := models.Patient{
		Name:               name,
		Id:                 id,
		Address:            address,
		PhoneNumber:        phoneNumber,
		ExternalIdentifier: externalIdentifier,
		DateOfBirth:        dateOfBirth,
	}

	err = s.repo.UpdatePatient(ctx, patient)
	if err != nil {
		return models.Patient{}, fmt.Errorf("error updating patient %w", err)
	}
	return patient, nil
}

func (s PatientService) DeletePatient(ctx context.Context, patientId int) error {
	ctx, span := s.tracer.NewSpan(ctx, "DeletePatient")
	defer span.End()
	s.tracer.SetAttributes(ctx, attribute.Int("patientId", patientId))

	err := s.ValidatePatientId(ctx, patientId)
	if err != nil {
		return err
	}

	err = s.repo.DeletePatient(ctx, patientId)
	if err != nil {
		return s.tracer.RecordError(ctx, fmt.Errorf("error deleting patient %w", err))
	}

	return nil
}

func (s PatientService) SearchPatients(ctx context.Context, search models.PatientSearch) ([]models.Patient, error) {
	ctx, span := s.tracer.NewSpan(ctx, "SearchPatients")
	defer span.End()
	s.tracer.SetAttributes(ctx,
		attribute.String("search.address", search.Address),
		attribute.String("search.attatchmentName", search.AttatchmentName),
		attribute.String("search.attatchmentType", search.AttatchmentType),
		attribute.String("search.diagnosedConditionCode", search.DiagnosedConditionCode),
		attribute.String("search.diagnosedConditionName", search.DiagnosedConditionName),
		attribute.String("search.name", search.Name),
		attribute.String("search.phone", search.Phone),
	)

	patients, err := s.repo.SearchPatients(ctx, search)
	if err != nil {
		return nil, s.tracer.RecordError(ctx, fmt.Errorf("error searching patients %w", err))
	}

	return patients, nil
}

func (s PatientService) ValidatePatientId(ctx context.Context, patientId int) error {
	count, err := s.repo.GetCountOfPatientId(ctx, patientId)
	if err != nil {
		return s.tracer.RecordError(ctx, fmt.Errorf("error getting patient id count %w", err))
	}

	if count < 1 {
		return s.tracer.RecordError(ctx, customerrors.NewInvalidInputError("patient id not found"))
	}
	return nil
}
