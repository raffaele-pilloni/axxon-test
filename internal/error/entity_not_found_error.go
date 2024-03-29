package error

import "fmt"

const (
	MessageEntityNotFoundErrorMessage string = "There is no %s with id %d."
)

type EntityNotFoundError struct {
	code    int
	message string
}

func NewEntityNotFoundError(
	entity string,
	entityID int,
) *EntityNotFoundError {
	return &EntityNotFoundError{
		code:    EntityNotFoundErrorCode,
		message: fmt.Sprintf(MessageEntityNotFoundErrorMessage, entity, entityID),
	}
}

func (e *EntityNotFoundError) Code() int {
	return e.code
}

func (e *EntityNotFoundError) Message() string {
	return e.message
}

func (e *EntityNotFoundError) Error() string {
	return fmt.Sprintf(
		"Code %d: Message: %s",
		e.code,
		e.message,
	)
}
