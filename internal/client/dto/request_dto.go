package dto

type RequestDTO struct {
	Method  string
	URL     string
	Headers map[string][]string
	Body    []byte
}
