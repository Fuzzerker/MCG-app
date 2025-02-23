package diagnosedconditions

import (
	"context"
	"fmt"
	"mcg-app-backend/service/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MockDiagnosedConditionRepo struct {
	mock.Mock
}

func (m *MockDiagnosedConditionRepo) InsertDiagnosedCondition(ctx context.Context, condition models.DiagnosedCondition) (int, error) {
	args := m.Called(ctx, condition)
	return args.Int(0), args.Error(1)
}

func (m *MockDiagnosedConditionRepo) DeleteDiagnosedCondition(ctx context.Context, conditionId int) error {
	args := m.Called(ctx, conditionId)
	return args.Error(0)
}

func (m *MockDiagnosedConditionRepo) DeleteDiagnosedConditionsByPatientId(ctx context.Context, patientId int) error {
	args := m.Called(ctx, patientId)
	return args.Error(0)
}

type MockPatientService struct {
	mock.Mock
}

func (m *MockPatientService) ValidatePatientId(ctx context.Context, patientId int) error {
	args := m.Called(ctx, patientId)
	return args.Error(0)
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

func getMocksAndService() (*MockDiagnosedConditionRepo, *MockPatientService, DiagnosedConditionService) {
	mockRepo := new(MockDiagnosedConditionRepo)
	mockPatientSvc := new(MockPatientService)
	mockTracer := new(MockTracer)
	service := NewDiagnosedConditionService(mockRepo, mockPatientSvc, mockTracer)
	return mockRepo, mockPatientSvc, service
}

func TestAddDiagnosedConditionToPatient(t *testing.T) {
	patientId := 1
	name := "Diabetes"
	code := "ABC123"
	description := "Type 2 Diabetes"
	date := time.Now()

	t.Run("AddDiagnosedCondition_Success", func(t *testing.T) {
		mockRepo, mockPatientSvc, service := getMocksAndService()
		mockPatientSvc.On("ValidatePatientId", mock.Anything, patientId).Return(nil)
		mockRepo.On("InsertDiagnosedCondition", mock.Anything, mock.AnythingOfType("models.DiagnosedCondition")).Return(1, nil)

		condition, err := service.AddDiagnosedConditionToPatient(context.Background(), patientId, name, code, description, date)
		assert.Nil(t, err)
		assert.Equal(t, name, condition.Name)
		assert.Equal(t, patientId, condition.PatientId)

		mockRepo.AssertExpectations(t)
		mockPatientSvc.AssertExpectations(t)
	})

	t.Run("AddDiagnosedCondition_InvalidPatientId", func(t *testing.T) {
		mockRepo, mockPatientSvc, service := getMocksAndService()
		mockPatientSvc.On("ValidatePatientId", mock.Anything, patientId).Return(fmt.Errorf("patient id not found"))

		condition, err := service.AddDiagnosedConditionToPatient(context.Background(), patientId, name, code, description, date)
		assert.NotNil(t, err)
		assert.Empty(t, condition)
		assert.Equal(t, "patient id not found", err.Error())

		mockRepo.AssertExpectations(t)
		mockPatientSvc.AssertExpectations(t)
	})

	t.Run("AddDiagnosedCondition_RepoError", func(t *testing.T) {
		mockRepo, mockPatientSvc, service := getMocksAndService()
		mockPatientSvc.On("ValidatePatientId", mock.Anything, patientId).Return(nil)
		mockRepo.On("InsertDiagnosedCondition", mock.Anything, mock.AnythingOfType("models.DiagnosedCondition")).Return(0, fmt.Errorf("db error"))

		condition, err := service.AddDiagnosedConditionToPatient(context.Background(), patientId, name, code, description, date)
		assert.NotNil(t, err)
		assert.Empty(t, condition)
		assert.Equal(t, "error inserting diagnosed condition db error", err.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteDiagnosedCondition(t *testing.T) {
	conditionId := 1

	t.Run("DeleteDiagnosedCondition_Success", func(t *testing.T) {
		mockRepo, _, service := getMocksAndService()
		mockRepo.On("DeleteDiagnosedCondition", mock.Anything, conditionId).Return(nil)

		err := service.DeleteDiagnosedCondition(context.Background(), conditionId)
		assert.Nil(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteDiagnosedCondition_RepoError", func(t *testing.T) {
		mockRepo, _, service := getMocksAndService()
		mockRepo.On("DeleteDiagnosedCondition", mock.Anything, conditionId).Return(fmt.Errorf("db error"))

		err := service.DeleteDiagnosedCondition(context.Background(), conditionId)
		assert.NotNil(t, err)
		assert.Equal(t, "error deleting diagnosed condition db error", err.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestDeletePatientDiagnosedConditions(t *testing.T) {
	patientId := 1

	t.Run("DeletePatientDiagnosedConditions_Success", func(t *testing.T) {
		mockRepo, _, service := getMocksAndService()
		mockRepo.On("DeleteDiagnosedConditionsByPatientId", mock.Anything, patientId).Return(nil)

		err := service.DeletePatientDiagnosedConditions(context.Background(), patientId)
		assert.Nil(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeletePatientDiagnosedConditions_RepoError", func(t *testing.T) {
		mockRepo, _, service := getMocksAndService()
		mockRepo.On("DeleteDiagnosedConditionsByPatientId", mock.Anything, patientId).Return(fmt.Errorf("db error"))

		err := service.DeletePatientDiagnosedConditions(context.Background(), patientId)
		assert.NotNil(t, err)
		assert.Equal(t, "error deleting diagnosed conditions for patient db error", err.Error())

		mockRepo.AssertExpectations(t)
	})
}
