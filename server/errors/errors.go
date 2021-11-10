package errors

type Error struct {
	Message string
	Inner   error
}

func BadRequest(inner error) *Error { return &Error{"Bad request", inner} }
func FileNotFound() *Error          { return &Error{"File not found", nil} }

func (e *Error) Error() string {
	return e.Message
}
