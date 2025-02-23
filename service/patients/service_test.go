package patients

import (
	"context"
	"errors"
	"fmt"
	"mcg-app-backend/service/customerrors"
	"mcg-app-backend/service/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MockPatientRepo struct {
	mock.Mock
}

func (m *MockPatientRepo) GetCountOfExternalIdentifier(ctx context.Context, externalIdentifier string) (int, error) {
	args := m.Called(ctx, externalIdentifier)
	return args.Int(0), args.Error(1)
}

func (m *MockPatientRepo) InsertPatient(ctx context.Context, patient models.Patient) (int, error) {
	args := m.Called(ctx, patient)
	return args.Int(0), args.Error(1)
}

func (m *MockPatientRepo) UpdatePatient(ctx context.Context, patient models.Patient) error {
	args := m.Called(ctx, patient)
	return args.Error(0)
}

func (m *MockPatientRepo) DeletePatient(ctx context.Context, patientId int) error {
	args := m.Called(ctx, patientId)
	return args.Error(0)
}

func (m *MockPatientRepo) SearchPatients(ctx context.Context, search models.PatientSearch) ([]models.Patient, error) {
	args := m.Called(ctx, search)
	return args.Get(0).([]models.Patient), args.Error(1)
}

func (m *MockPatientRepo) GetCountOfPatientId(ctx context.Context, patientId int) (int, error) {
	args := m.Called(ctx, patientId)
	return args.Int(0), args.Error(1)
}

type MockTracer struct {
	mock.Mock
}

func (m *MockTracer) NewSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer("").Start(context.Background(), "")
}

func (m *MockTracer) SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {}

func (m *MockTracer) RecordError(ctx context.Context, err error) error {
	return err
}

func getMocksAndService() (*MockPatientRepo, PatientService) {
	mockRepo := new(MockPatientRepo)
	mockTracer := new(MockTracer)
	service := NewPatientService(mockRepo, mockTracer)
	return mockRepo, service
}

func TestCreatePatient(t *testing.T) {
	name := "John Doe"
	address := "123 Main St"
	phoneNumber := "1234567890"
	externalIdentifier := "unique123"

	t.Run("CreatePatient_Success", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetCountOfExternalIdentifier", mock.Anything, externalIdentifier).Return(0, nil)
		mockRepo.On("InsertPatient", mock.Anything, mock.AnythingOfType("models.Patient")).Return(1, nil)

		patient, err := service.CreatePatient(context.Background(), name, address, phoneNumber, time.Now(), externalIdentifier)
		assert.Nil(t, err)
		assert.Equal(t, name, patient.Name)

		mockRepo.AssertExpectations(t)
	})

	t.Run("CreatePatient_ExternalIdentifierAlreadyExists", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetCountOfExternalIdentifier", mock.Anything, externalIdentifier).Return(1, nil)

		patient, err := service.CreatePatient(context.Background(), name, address, phoneNumber, time.Now(), externalIdentifier)
		assert.NotNil(t, err)
		assert.Empty(t, patient)
		var alreadyExists customerrors.AlreadyExistsError
		assert.True(t, errors.As(err, &alreadyExists))

		mockRepo.AssertExpectations(t)
	})

	t.Run("CreatePatient_RepoError_GetCountOfExternalIdentifier", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetCountOfExternalIdentifier", mock.Anything, externalIdentifier).Return(0, fmt.Errorf("db error"))

		patient, err := service.CreatePatient(context.Background(), name, address, phoneNumber, time.Now(), externalIdentifier)
		assert.NotNil(t, err)
		assert.Empty(t, patient)

		mockRepo.AssertExpectations(t)
	})

	t.Run("CreatePatient_RepoError_InsertPatient", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetCountOfExternalIdentifier", mock.Anything, externalIdentifier).Return(0, nil)
		mockRepo.On("InsertPatient", mock.Anything, mock.AnythingOfType("models.Patient")).Return(0, fmt.Errorf("db error"))

		patient, err := service.CreatePatient(context.Background(), name, address, phoneNumber, time.Now(), externalIdentifier)
		assert.NotNil(t, err)
		assert.Empty(t, patient)

		mockRepo.AssertExpectations(t)
	})
}

func TestUpdatePatient(t *testing.T) {
	patient := models.Patient{
		Id:          1,
		Name:        "Jane Doe",
		PhoneNumber: "9876543210",
		Address:     "456 Elm St",
	}

	t.Run("UpdatePatient_Success", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetCountOfPatientId", mock.Anything, patient.Id).Return(1, nil) // Mocking the patient ID count to be 1 (valid)
		mockRepo.On("UpdatePatient", mock.Anything, patient).Return(nil)

		updatedPatient, err := service.UpdatePatient(context.Background(), patient.Id, patient.Name, patient.Address, patient.PhoneNumber, patient.DateOfBirth, patient.ExternalIdentifier)
		assert.Nil(t, err)
		assert.Equal(t, patient.Name, updatedPatient.Name)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdatePatient_InvalidPatientId", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetCountOfPatientId", mock.Anything, patient.Id).Return(0, nil) // Mocking the patient ID count to be 0 (invalid)

		updatedPatient, err := service.UpdatePatient(context.Background(), patient.Id, patient.Name, patient.Address, patient.PhoneNumber, patient.DateOfBirth, patient.ExternalIdentifier)
		assert.NotNil(t, err)
		assert.Equal(t, "patient id not found", err.Error())
		assert.Empty(t, updatedPatient)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdatePatient_RepoError", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetCountOfPatientId", mock.Anything, patient.Id).Return(1, nil)
		mockRepo.On("UpdatePatient", mock.Anything, patient).Return(fmt.Errorf("db error"))

		updatedPatient, err := service.UpdatePatient(context.Background(), patient.Id, patient.Name, patient.Address, patient.PhoneNumber, patient.DateOfBirth, patient.ExternalIdentifier)
		assert.NotNil(t, err)
		assert.Equal(t, "error updating patient db error", err.Error())
		assert.Empty(t, updatedPatient)

		mockRepo.AssertExpectations(t)
	})
}

func TestDeletePatient(t *testing.T) {
	patientId := 1

	t.Run("DeletePatient_Success", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetCountOfPatientId", mock.Anything, patientId).Return(1, nil) // Mocking the patient ID count to be 1 (valid)
		mockRepo.On("DeletePatient", mock.Anything, patientId).Return(nil)

		err := service.DeletePatient(context.Background(), patientId)
		assert.Nil(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeletePatient_InvalidPatientId", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetCountOfPatientId", mock.Anything, patientId).Return(0, nil) // Mocking the patient ID count to be 0 (invalid)

		err := service.DeletePatient(context.Background(), patientId)
		assert.NotNil(t, err)
		assert.Equal(t, "patient id not found", err.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeletePatient_RepoError", func(t *testing.T) {
		mockRepo, service := getMocksAndService()
		mockRepo.On("GetCountOfPatientId", mock.Anything, patientId).Return(1, nil)
		mockRepo.On("DeletePatient", mock.Anything, patientId).Return(fmt.Errorf("db error"))

		err := service.DeletePatient(context.Background(), patientId)
		assert.NotNil(t, err)
		assert.Equal(t, "error deleting patient db error", err.Error())

		mockRepo.AssertExpectations(t)
	})
}
