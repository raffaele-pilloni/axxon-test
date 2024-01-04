package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	httperror "github.com/raffaele-pilloni/axxon-test/internal/app/http/error"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/model/response"
	applicationerror "github.com/raffaele-pilloni/axxon-test/internal/error"
	"net/http"
)

const (
	ErrorCodeInternalServerError = 0
	ErrorCodeEntityNotFoundError = 1
	ErrorCodeBadRequest          = 2

	ErrorMessageInternalServerError = "something went wrong"
)

type successModelResponse interface {
	response.CreateTaskModelResponse |
		response.GetTaskModelResponse
}

func HandleSuccess[T successModelResponse](w http.ResponseWriter, successModelResponse *T) {
	successModelResponseMarshalled, err := json.Marshal(successModelResponse)
	if err != nil {
		HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(successModelResponseMarshalled)
}

func HandleError(w http.ResponseWriter, error error) {
	switch error.(type) {
	case applicationerror.ApplicationErrorInterface:
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorModelResponseMarshalled, err := json.Marshal(&response.ErrorModelResponse{
			Code:    error.(applicationerror.ApplicationErrorInterface).Code(),
			Message: error.(applicationerror.ApplicationErrorInterface).Message(),
		})
		if err != nil {
			return
		}

		w.Write(errorModelResponseMarshalled)
	case *httperror.InvalidJSONError, *validator.InvalidValidationError, validator.ValidationErrors:
		w.WriteHeader(http.StatusBadRequest)
		errorModelResponseMarshalled, err := json.Marshal(&response.ErrorModelResponse{
			Code:    ErrorCodeBadRequest,
			Message: error.Error(),
		})
		if err != nil {
			return
		}

		w.Write(errorModelResponseMarshalled)
	case *httperror.EntityNotFoundError:
		w.WriteHeader(http.StatusNotFound)
		errorModelResponseMarshalled, err := json.Marshal(&response.ErrorModelResponse{
			Code:    ErrorCodeEntityNotFoundError,
			Message: error.Error(),
		})
		if err != nil {
			return
		}

		w.Write(errorModelResponseMarshalled)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		errorModelResponseMarshalled, err := json.Marshal(&response.ErrorModelResponse{
			Code:    ErrorCodeInternalServerError,
			Message: ErrorMessageInternalServerError,
		})
		if err != nil {
			return
		}

		w.Write(errorModelResponseMarshalled)
	}
}
