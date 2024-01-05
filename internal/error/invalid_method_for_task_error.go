package error

import (
	"fmt"
)

const (
	InvalidMethodForTaskErrorMessage = "The method %s is not valid for task."
)

type InvalidMethodForTaskError struct {
	code    int
	message string
}

func NewInvalidMethodForTaskError(method string) *InvalidMethodForTaskError {
	return &InvalidMethodForTaskError{
		code:    CodeInvalidMethodForTaskError,
		message: fmt.Sprintf(InvalidMethodForTaskErrorMessage, method),
	}
}

func (i *InvalidMethodForTaskError) Code() int {
	return i.code
}

func (i *InvalidMethodForTaskError) Message() string {
	return i.message
}

func (i *InvalidMethodForTaskError) Error() string {
	return fmt.Sprintf(
		"Code %d: Message: %s",
		i.code,
		i.message,
	)
}
