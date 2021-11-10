package errors

type Error struct {
	Message string `json:"error"`
	Cause   error  `json:"cause"`
}

func BadRequest(inner error) *Error { return &Error{"Bad request", inner} }
func FileNotFound() *Error          { return &Error{"File not found", nil} }

func (e *Error) Error() string {
	return e.Message
}
