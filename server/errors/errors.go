package errors

type Error struct {
	Message string `json:"error"`
	Cause   error  `json:"cause,omitempty"`
}

func BadRequest(inner error) *Error    { return &Error{"Bad request", inner} }
func FileNotFound() *Error             { return &Error{"File not found", nil} }
func MissingParam(param string) *Error { return &Error{"Missing parameter " + param, nil} }

func (e *Error) Error() string {
	return e.Message
}
