package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	httperror "github.com/raffaele-pilloni/axxon-test/internal/app/http/error"
	"github.com/raffaele-pilloni/axxon-test/internal/app/http/model/request"
	"net/http"
)

type modelRequest interface {
	request.CreateTaskModelRequest
}

func HandleRequest[T modelRequest](r *http.Request) (*T, error) {
	var modelRequest T

	if err := json.NewDecoder(r.Body).Decode(&modelRequest); err != nil {
		return nil, httperror.NewInvalidJSONError(err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(&modelRequest); err != nil {
		return nil, err
	}

	return &modelRequest, nil
}
