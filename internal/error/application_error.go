package error

const (
	// List of all the possible error codes related to DAL.
	CodeEntityNotFoundError int = 1000
	// List of all the possible error codes related to Task.
	CodeInvalidMethodForTaskError int = 2000
	CodeInvalidURLForTaskError    int = 2001
)

type ApplicationErrorInterface interface {
	error
	Code() int
	Message() string
}
