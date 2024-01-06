package dto

type RequestDTO struct {
	Method  string
	URL     string
	Body    []byte
	Headers map[string][]string
}
