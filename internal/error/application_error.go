package error

const (
	// List of all the possible error codes related to Task.
	ErrorCodeInvalidMethodForTask = 1000
	ErrorCodeInvalidURLForTask    = 1001
)

type ApplicationErrorInterface interface {
	error
	Code() int
	Message() string
}
