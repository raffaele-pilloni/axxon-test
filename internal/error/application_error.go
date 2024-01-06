package error

const (
	// List of all the possible error codes related to DAL.
	EntityNotFoundErrorCode int = 1000
	// List of all the possible error codes related to Task.
	InvalidMethodForTaskErrorCode int = 2000
	InvalidURLForTaskErrorCode    int = 2001
)

type ApplicationErrorInterface interface {
	error
	Code() int
	Message() string
}
