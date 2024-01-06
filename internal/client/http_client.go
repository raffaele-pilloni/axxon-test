package client

import (
	"bytes"
	"context"
	"github.com/raffaele-pilloni/axxon-test/internal/client/dto"
	"net/http"
)

type HTTPClientInterface interface {
	Do(ctx context.Context, requestDTO *dto.RequestDTO) (*dto.ResponseDTO, error)
}

type HTTPClient struct {
	httpClient *http.Client
}

func NewHTTPClient(
	httpClient *http.Client,
) *HTTPClient {
	return &HTTPClient{
		httpClient: httpClient,
	}
}

func (h *HTTPClient) Do(ctx context.Context, requestDTO *dto.RequestDTO) (*dto.ResponseDTO, error) {
	request, err := http.NewRequestWithContext(ctx, requestDTO.Method, requestDTO.URL, bytes.NewReader(requestDTO.Body))
	if err != nil {
		return nil, err
	}

	request.Header = requestDTO.Headers

	response, err := h.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return &dto.ResponseDTO{
		Header:        response.Header,
		StatusCode:    response.StatusCode,
		ContentLength: int(response.ContentLength),
	}, err
}
