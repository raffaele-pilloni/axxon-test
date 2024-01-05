package error

type InvalidRequestError struct {
	error error
}

func NewInvalidRequestError(
	error error,
) *InvalidRequestError {
	return &InvalidRequestError{
		error: error,
	}
}

func (i *InvalidRequestError) Error() string {
	return i.error.Error()
}
