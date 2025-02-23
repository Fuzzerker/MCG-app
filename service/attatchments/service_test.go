package attatchments

import (
	"context"
	"fmt"
	"mcg-app-backend/service/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MockAttachmentRepo struct {
	mock.Mock
}

func (m *MockAttachmentRepo) InsertAttatchment(ctx context.Context, attachment models.Attatchment) (int, error) {
	args := m.Called(ctx, attachment)
	return args.Int(0), args.Error(1)
}

func (m *MockAttachmentRepo) DeleteAttatchment(ctx context.Context, attachmentId int) error {
	args := m.Called(ctx, attachmentId)
	return args.Error(0)
}

func (m *MockAttachmentRepo) DeleteAttatchmentsByPatientId(ctx context.Context, patientId int) error {
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
func getMocksAndService() (*MockAttachmentRepo, *MockPatientService, *MockTracer, AttachmentService) {
	mockAttachmentRepo := new(MockAttachmentRepo)
	mockPatientService := new(MockPatientService)
	mockTracer := new(MockTracer)
	service := NewAttachmentService(mockAttachmentRepo, mockPatientService, mockTracer)
	return mockAttachmentRepo, mockPatientService, mockTracer, service
}

func TestAddAttatchmentToPatient(t *testing.T) {
	patientId := 1
	name := "Attachment 1"
	description := "Description of attachment"
	typ := "pdf"
	data := []byte("sample data")

	t.Run("AddAttachment_Success", func(t *testing.T) {
		mockAttachmentRepo, mockPatientService, mockTracer, service := getMocksAndService()
		mockPatientService.On("ValidatePatientId", mock.Anything, patientId).Return(nil)
		mockAttachmentRepo.On("InsertAttatchment", mock.Anything, mock.Anything).Return(1, nil)

		attachment, err := service.AddAttatchmentToPatient(context.Background(), patientId, name, description, typ, data)
		assert.Nil(t, err)
		assert.NotNil(t, attachment)
		assert.Equal(t, patientId, attachment.PatientId)
		assert.Equal(t, name, attachment.Name)
		assert.Equal(t, description, attachment.Description)
		assert.Equal(t, typ, attachment.Type)
		assert.Equal(t, data, attachment.Data)

		mockPatientService.AssertExpectations(t)
		mockAttachmentRepo.AssertExpectations(t)
		mockTracer.AssertExpectations(t)
	})

	t.Run("AddAttachment_EmptyData", func(t *testing.T) {
		mockAttachmentRepo, mockPatientService, mockTracer, service := getMocksAndService()

		attachment, err := service.AddAttatchmentToPatient(context.Background(), patientId, name, description, typ, []byte{})
		assert.NotNil(t, err)
		assert.Equal(t, "data was empty", err.Error())
		assert.Equal(t, models.Attatchment{}, attachment)

		mockPatientService.AssertExpectations(t)
		mockAttachmentRepo.AssertExpectations(t)
		mockTracer.AssertExpectations(t)
	})

	t.Run("AddAttachment_InvalidPatient", func(t *testing.T) {
		mockAttachmentRepo, mockPatientService, mockTracer, service := getMocksAndService()
		mockPatientService.On("ValidatePatientId", mock.Anything, patientId).Return(fmt.Errorf("invalid patient"))

		attachment, err := service.AddAttatchmentToPatient(context.Background(), patientId, name, description, typ, data)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid patient", err.Error())
		assert.Equal(t, models.Attatchment{}, attachment)

		mockPatientService.AssertExpectations(t)
		mockAttachmentRepo.AssertExpectations(t)
		mockTracer.AssertExpectations(t)
	})

	t.Run("AddAttachment_InsertError", func(t *testing.T) {
		mockAttachmentRepo, mockPatientService, mockTracer, service := getMocksAndService()
		mockPatientService.On("ValidatePatientId", mock.Anything, patientId).Return(nil)
		mockAttachmentRepo.On("InsertAttatchment", mock.Anything, mock.Anything).Return(0, fmt.Errorf("insert error"))

		attachment, err := service.AddAttatchmentToPatient(context.Background(), patientId, name, description, typ, data)
		assert.NotNil(t, err)
		assert.Equal(t, "error inserting attachment insert error", err.Error())
		assert.Equal(t, models.Attatchment{}, attachment)

		mockPatientService.AssertExpectations(t)
		mockAttachmentRepo.AssertExpectations(t)
		mockTracer.AssertExpectations(t)
	})
}

func TestDeleteAttatchment(t *testing.T) {
	attachmentId := 1

	t.Run("DeleteAttachment_Success", func(t *testing.T) {
		mockAttachmentRepo, _, mockTracer, service := getMocksAndService()
		mockAttachmentRepo.On("DeleteAttatchment", mock.Anything, attachmentId).Return(nil)

		err := service.DeleteAttatchment(context.Background(), attachmentId)
		assert.Nil(t, err)

		mockAttachmentRepo.AssertExpectations(t)
		mockTracer.AssertExpectations(t)
	})

	t.Run("DeleteAttachment_Error", func(t *testing.T) {
		mockAttachmentRepo, _, mockTracer, service := getMocksAndService()
		mockAttachmentRepo.On("DeleteAttatchment", mock.Anything, attachmentId).Return(fmt.Errorf("delete error"))

		err := service.DeleteAttatchment(context.Background(), attachmentId)
		assert.NotNil(t, err)
		assert.Equal(t, "error deleting attachment delete error", err.Error())

		mockAttachmentRepo.AssertExpectations(t)
		mockTracer.AssertExpectations(t)
	})
}

func TestDeletePatientAttachments(t *testing.T) {
	patientId := 1

	t.Run("DeletePatientAttachments_Success", func(t *testing.T) {
		mockAttachmentRepo, _, mockTracer, service := getMocksAndService()
		mockAttachmentRepo.On("DeleteAttatchmentsByPatientId", mock.Anything, patientId).Return(nil)

		err := service.DeletePatientAttachments(context.Background(), patientId)
		assert.Nil(t, err)

		mockAttachmentRepo.AssertExpectations(t)
		mockTracer.AssertExpectations(t)
	})

	t.Run("DeletePatientAttachments_Error", func(t *testing.T) {
		mockAttachmentRepo, _, mockTracer, service := getMocksAndService()
		mockAttachmentRepo.On("DeleteAttatchmentsByPatientId", mock.Anything, patientId).Return(fmt.Errorf("delete error"))

		err := service.DeletePatientAttachments(context.Background(), patientId)
		assert.NotNil(t, err)
		assert.Equal(t, "error deleting attachments for patient delete error", err.Error())

		mockAttachmentRepo.AssertExpectations(t)
		mockTracer.AssertExpectations(t)
	})
}
