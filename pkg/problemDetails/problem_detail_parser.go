package problemDetails

import (
	"context"
	"database/sql"
	"github.com/go-playground/validator"
	customErrors2 "github.com/meysamhadeli/shop-golang-microservices/pkg/problemDetails/custome_error"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/problemDetails/custome_error/utils"
	"github.com/pkg/errors"
	"net/http"
)

func ParseError(err error) ProblemDetailErr {
	stackTrace := httpErrors.ErrorsWithStack(err)
	customErr := customErrors2.GetCustomError(err)
	var validatorErr validator.ValidationErrors

	if err != nil {
		switch {
		case customErrors2.IsDomainError(err):
			return NewDomainProblemDetail(customErr.Status(), customErr.Error(), stackTrace)
		case customErrors2.IsApplicationError(err):
			return NewApplicationProblemDetail(customErr.Status(), customErr.Error(), stackTrace)
		case customErrors2.IsApiError(err):
			return NewApiProblemDetail(customErr.Status(), customErr.Error(), stackTrace)
		case customErrors2.IsBadRequestError(err):
			return NewBadRequestProblemDetail(customErr.Error(), stackTrace)
		case customErrors2.IsNotFoundError(err):
			return NewNotFoundErrorProblemDetail(customErr.Error(), stackTrace)
		case customErrors2.IsValidationError(err):
			return NewValidationProblemDetail(customErr.Error(), stackTrace)
		case customErrors2.IsUnAuthorizedError(err):
			return NewUnAuthorizedErrorProblemDetail(customErr.Error(), stackTrace)
		case customErrors2.IsForbiddenError(err):
			return NewForbiddenProblemDetail(customErr.Error(), stackTrace)
		case customErrors2.IsConflictError(err):
			return NewConflictProblemDetail(customErr.Error(), stackTrace)
		case customErrors2.IsInternalServerError(err):
			return NewInternalServerProblemDetail(customErr.Error(), stackTrace)
		case customErrors2.IsCustomError(err):
			return NewProblemDetailFromCodeAndDetail(customErr.Status(), customErr.Error(), stackTrace)
		case customErrors2.IsUnMarshalingError(err):
			return NewInternalServerProblemDetail(err.Error(), stackTrace)
		case customErrors2.IsMarshalingError(err):
			return NewInternalServerProblemDetail(err.Error(), stackTrace)
		case errors.Is(err, sql.ErrNoRows):
			return NewNotFoundErrorProblemDetail(err.Error(), stackTrace)
		case errors.Is(err, context.DeadlineExceeded):
			return NewProblemDetail(http.StatusRequestTimeout, ErrRequestTimeoutTitle, err.Error(), stackTrace)
		case errors.As(err, &validatorErr):
			return NewValidationProblemDetail(validatorErr.Error(), stackTrace)
		default:
			return NewInternalServerProblemDetail(err.Error(), stackTrace)
		}
	}

	return nil
}
