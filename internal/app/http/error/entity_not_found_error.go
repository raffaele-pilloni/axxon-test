package error

import "fmt"

const (
	EntityNotFoundErrorMessage string = "There is no %s with id %d."
)

type EntityNotFoundError struct {
	message string
}

func NewEntityNotFoundError(
	entity string,
	entityID int,
) *EntityNotFoundError {
	return &EntityNotFoundError{
		message: fmt.Sprintf(EntityNotFoundErrorMessage, entity, entityID),
	}
}

func (e *EntityNotFoundError) Error() string {
	return e.message
}
