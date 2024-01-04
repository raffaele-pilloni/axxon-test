package error

import (
	"fmt"
)

type InvalidURLForTaskError struct {
	code    int
	message string
}

func NewInvalidURLForTaskError(url string) *InvalidURLForTaskError {
	return &InvalidURLForTaskError{
		code:    ErrorCodeInvalidURLForTask,
		message: fmt.Sprintf("The url %s is not valid for task.", url),
	}
}

func (i *InvalidURLForTaskError) Code() int {
	return i.code
}

func (i *InvalidURLForTaskError) Message() string {
	return i.message
}

func (i *InvalidURLForTaskError) Error() string {
	return fmt.Sprintf(
		"Code %d: Message: %s",
		i.code,
		i.message,
	)
}
