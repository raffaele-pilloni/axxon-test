package dto

type ResponseDTO struct {
	Headers    map[string][]string
	StatusCode int
	Body       []byte
}
