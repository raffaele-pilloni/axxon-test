package dto

type ResponseDTO struct {
	Header     map[string][]string
	StatusCode int
	Body       []byte
}
