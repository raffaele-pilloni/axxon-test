package dto

type DoneTaskProcessingDTO struct {
	ResponseHeaders       map[string][]string
	ResponseStatusCode    int
	ResponseContentLength int
}
