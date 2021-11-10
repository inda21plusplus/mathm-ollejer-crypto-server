package server

type Error struct {
	message string
	inner   error
}

func BadRequest(inner error) *Error { return &Error{"Bad request", inner} }

func (e *Error) Error() string {
	return e.message
}
