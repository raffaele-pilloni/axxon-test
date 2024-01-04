package error

type InvalidJSONError struct {
	error error
}

func NewInvalidJSONError(
	error error,
) *InvalidJSONError {
	return &InvalidJSONError{
		error: error,
	}
}

func (i *InvalidJSONError) Error() string {
	return i.error.Error()
}
