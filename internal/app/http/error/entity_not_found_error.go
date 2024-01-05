package error

import "fmt"

const (
	ErrorMessageEntityNotFoundError string = "There is no %s with %d."
)

type EntityNotFoundError struct {
	message string
}

func NewEntityNotFoundError(
	entity string,
	entityID int,
) *EntityNotFoundError {
	return &EntityNotFoundError{
		message: fmt.Sprintf(ErrorMessageEntityNotFoundError, entity, entityID),
	}
}

func (e *EntityNotFoundError) Error() string {
	return e.message
}
