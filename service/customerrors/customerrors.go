package customerrors

import "net/http"

type AlreadyExistsError struct {
	message string
}

func (r AlreadyExistsError) Error() string {
	return r.message
}

func (r AlreadyExistsError) HTTPStatus() int {
	return http.StatusConflict
}

func NewAlreadyExistsError(message string) AlreadyExistsError {
	return AlreadyExistsError{
		message: message,
	}
}

type InvalidInputError struct {
	message string
}

func (r InvalidInputError) Error() string {
	return r.message
}

func (r InvalidInputError) HTTPStatus() int {
	return http.StatusBadRequest
}

func NewInvalidInputError(message string) InvalidInputError {
	return InvalidInputError{
		message: message,
	}
}

type UnauthorizedError struct {
	message string
}

func (r UnauthorizedError) Error() string {
	return r.message
}

func (r UnauthorizedError) HTTPStatus() int {
	return http.StatusBadRequest
}

func NewUnauthorizedError(message string) UnauthorizedError {
	return UnauthorizedError{
		message: message,
	}
}
