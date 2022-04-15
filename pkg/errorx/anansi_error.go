package errorx

type ErrorType string

const (
	NotFoundError     ErrorType = "Not Found"
	ValidationFailure ErrorType = "Validation Failure"
)

type AnansiError struct {
	Error   error
	Type    ErrorType
	Message string
}

func NewError(originalError error, errorType ErrorType, message string) *AnansiError {
	return &AnansiError{
		Error:   originalError,
		Type:    errorType,
		Message: message,
	}
}

func (e *AnansiError) GetError() error {
	return e.Error
}

func (e *AnansiError) IsNotFound() bool {
	return e.Type == NotFoundError
}

func (e *AnansiError) IsValidationFailure() bool {
	return e.Type == ValidationFailure
}
