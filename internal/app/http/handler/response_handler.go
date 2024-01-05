package handler

import (
	"encoding/json"
	httperror "github.com/raffaele-pilloni/axxon-test/internal/app/http/error"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/model/response"
	applicationerror "github.com/raffaele-pilloni/axxon-test/internal/error"
	"net/http"
)

const (
	// List of all the possible error codes related to DAL.
	HTTPInternalServerErrorCode = 0
	HTTPEntityNotFoundErrorCode = 1
	HTTPBadRequestErrorCode     = 2

	InternalServerErrorMessage = "Something went wrong."
)

type successModelResponse interface {
	response.CreateTaskModelResponse | response.GetTaskModelResponse
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
	case *httperror.EntityNotFoundError:
		w.WriteHeader(http.StatusNotFound)
		errorModelResponseMarshalled, err := json.Marshal(&response.ErrorModelResponse{
			Code:    HTTPEntityNotFoundErrorCode,
			Message: error.Error(),
		})
		if err != nil {
			return
		}

		w.Write(errorModelResponseMarshalled)
	case *httperror.InvalidRequestError:
		w.WriteHeader(http.StatusBadRequest)
		errorModelResponseMarshalled, err := json.Marshal(&response.ErrorModelResponse{
			Code:    HTTPBadRequestErrorCode,
			Message: error.Error(),
		})
		if err != nil {
			return
		}

		w.Write(errorModelResponseMarshalled)
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
	default:
		w.WriteHeader(http.StatusInternalServerError)
		errorModelResponseMarshalled, err := json.Marshal(&response.ErrorModelResponse{
			Code:    HTTPInternalServerErrorCode,
			Message: InternalServerErrorMessage,
		})
		if err != nil {
			return
		}

		w.Write(errorModelResponseMarshalled)
	}
}
