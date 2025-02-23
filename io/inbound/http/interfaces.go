package inboundhttp

import (
	"context"
	"mcg-app-backend/service/models"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, username string, password string) error
}

type AuthService interface {
	VerifyToken(ctx context.Context, tokenString string) error
	Login(ctx context.Context, username string, password string) (string, error)
}

type PatientService interface {
	CreatePatient(ctx context.Context, name string, address string, phoneNumber string, dateOfBirth time.Time, externalIdentifier string) (models.Patient, error)
	SearchPatients(ctx context.Context, search models.PatientSearch) ([]models.Patient, error)
	DeletePatient(ctx context.Context, patientId int) error
	UpdatePatient(ctx context.Context, id int, name string, address string, phoneNumber string, dateOfBirth time.Time, externalIdentifier string) (models.Patient, error)
}

type AttatchmentService interface {
	AddAttatchmentToPatient(ctx context.Context, patientId int, name string, description string, typ string, data []byte) (models.Attatchment, error)
	DeleteAttatchment(ctx context.Context, attatchmentId int) error
}

type DiagnosedConditionsService interface {
	AddDiagnosedConditionToPatient(ctx context.Context, patientId int, name string, code string, description string, date time.Time) (models.DiagnosedCondition, error)
	DeleteDiagnosedCondition(ctx context.Context, conditionId int) error
}
