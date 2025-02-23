package inboundhttp

import (
	"context"
	"errors"
	"io"
	"mcg-app-backend/service/customerrors"
	"mcg-app-backend/service/models"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (server HttpServer) handlePostUser() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input UserRequest, output *Empty) error {
		return handleError(server.userService.CreateUser(ctx, input.Username, input.Password))
	})
	u.SetExpectedErrors(status.InvalidArgument)
	u.SetExpectedErrors(status.AlreadyExists)
	return u
}

func (server HttpServer) handleLogin() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input UserRequest, output *LoginResponse) error {
		token, err := server.authService.Login(ctx, input.Username, input.Password)
		if err != nil {
			return handleError(err)
		}

		*output = LoginResponse{
			Token: token,
		}
		return nil

	})
	u.SetExpectedErrors(status.InvalidArgument)
	return u
}

func (server HttpServer) handleGetPatients() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input models.PatientSearch, output *[]models.Patient) error {
		patients, err := server.patientService.SearchPatients(ctx, input)
		if err != nil {
			return handleError(err)
		}

		*output = patients
		return nil

	})
	u.SetExpectedErrors(status.InvalidArgument)
	return u
}

func (server HttpServer) handleDeleteDiagnosedCondition() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input models.DeleteByIdRequest, output *Empty) error {
		err := server.diagnosedConditionService.DeleteDiagnosedCondition(ctx, input.Id)
		if err != nil {
			return handleError(err)
		}
		return nil

	})
	u.SetExpectedErrors(status.Unauthenticated)
	return u
}

func (server HttpServer) handleDeleteAttatchment() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input models.DeleteByIdRequest, output *Empty) error {
		err := server.attatchmentService.DeleteAttatchment(ctx, input.Id)
		if err != nil {
			return handleError(err)
		}
		return nil

	})
	u.SetExpectedErrors(status.Unauthenticated)
	return u
}

func (server HttpServer) handleDeletePatient() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input models.DeleteByIdRequest, output *Empty) error {
		err := server.patientService.DeletePatient(ctx, input.Id)
		if err != nil {
			return handleError(err)
		}
		return nil

	})
	u.SetExpectedErrors(status.Unauthenticated)
	return u
}

func (server HttpServer) handlePostPatient() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input models.Patient, output *models.Patient) error {
		patient, err := server.patientService.CreatePatient(ctx, input.Name, input.Address, input.PhoneNumber, input.ExternalIdentifier)
		if err != nil {
			return handleError(err)
		}

		*output = patient
		return nil

	})
	u.SetExpectedErrors(status.InvalidArgument)
	u.SetExpectedErrors(status.AlreadyExists)
	u.SetExpectedErrors(status.Unauthenticated)
	return u
}

func (server HttpServer) handlePutPatient() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input models.UpdatePatientRequest, output *models.Patient) error {
		input.Patient.Id = input.Id

		patient, err := server.patientService.UpdatePatient(ctx, input.Patient)
		if err != nil {
			return handleError(err)
		}

		*output = patient
		return nil

	})
	u.SetExpectedErrors(status.InvalidArgument)
	u.SetExpectedErrors(status.Unauthenticated)
	return u
}

func (server HttpServer) handlePostPatientAttatchment() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input models.CreateAttatchmentRequest, output *models.Attatchment) error {
		data, _ := io.ReadAll(input.Data)
		attatchment, err := server.attatchmentService.AddAttatchmentToPatient(ctx, input.PatientId, input.Name, input.Description, input.Type, data)

		if err != nil {
			return handleError(err)
		}

		*output = attatchment
		return nil

	})
	u.SetExpectedErrors(status.InvalidArgument)
	u.SetExpectedErrors(status.Unauthenticated)
	return u
}

func (server HttpServer) handlePostDiagnosedCondition() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input models.CreateDiagnosedConditionRequest, output *models.DiagnosedCondition) error {
		cond, err := server.diagnosedConditionService.AddDiagnosedConditionToPatient(ctx, input.PatientId, input.Name, input.Code, input.Description, input.Date)
		if err != nil {
			return handleError(err)
		}

		*output = cond
		return nil

	})
	u.SetExpectedErrors(status.InvalidArgument)
	u.SetExpectedErrors(status.Unauthenticated)
	return u
}

func handleError(err error) error {
	if err == nil {
		return nil
	}
	var inputError customerrors.InvalidInputError
	if errors.As(err, &inputError) {
		return inputError
	}
	var alreadyExistsError customerrors.AlreadyExistsError
	if errors.As(err, &alreadyExistsError) {
		return alreadyExistsError
	}
	var authError customerrors.UnauthorizedError
	if errors.As(err, &authError) {
		return authError
	}
	return errors.New("internal server error")
}
