package error

const (
	// List of all the possible error codes related to DAL.
	CodeEntityNotFoundError = 1000
	// List of all the possible error codes related to Task.
	CodeInvalidMethodForTaskError = 2000
	CodeInvalidURLForTaskError    = 2001
)

type ApplicationErrorInterface interface {
	error
	Code() int
	Message() string
}
