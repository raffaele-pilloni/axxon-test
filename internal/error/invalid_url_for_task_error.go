package error

import (
	"fmt"
)

const (
	InvalidURLForTaskErrorMessage string = "The url %s is not valid for task."
)

type InvalidURLForTaskError struct {
	code    int
	message string
}

func NewInvalidURLForTaskError(url string) *InvalidURLForTaskError {
	return &InvalidURLForTaskError{
		code:    CodeInvalidURLForTaskError,
		message: fmt.Sprintf(InvalidURLForTaskErrorMessage, url),
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
