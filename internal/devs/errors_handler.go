package devs

import (
	"errors"
	"net/http"

	"github.com/CriciumaDevJobs/backend/handlers"
)

var errorToStatus = map[error]int{
	ErrNameNotEmpty:         http.StatusUnprocessableEntity,
	ErrEmailNotEmpty:        http.StatusUnprocessableEntity,
	ErrEmailAddressNotValid: http.StatusUnprocessableEntity,
	ErrEmailAlreadyInUse:    http.StatusUnprocessableEntity,
	ErrBioNotEmpty:          http.StatusUnprocessableEntity,
	ErrPasswordToLong:       http.StatusUnprocessableEntity,
	ErrPasswordNotEmpty:     http.StatusUnprocessableEntity,

	ErrJsonNotExpected: http.StatusBadRequest,

	ErrInvalidEmailOrPassword: http.StatusUnauthorized,
	ErrProfileNotFound:        http.StatusNotFound,
}

func CheckUseCaseErr(err error) *handlers.ErrorResponse {

	status, ok := errorToStatus[err]

	if ok {
		return handlers.NewError(status, err.Error())
	}

	for sentinela, status := range errorToStatus {
		if errors.Is(err, sentinela) {
			return handlers.NewError(status, err.Error())
		}
	}

	return handlers.NewError(http.StatusInternalServerError, "Erro Interno!")
}
