package errors

type Error struct {
	Message string `json:"error"`
	Field   string `json:"field,omitempty"`
	Cause   error  `json:"cause,omitempty"`
}

func BadRequest(inner error) *Error    { return &Error{"bad request", "", inner} }
func FileNotFound() *Error             { return &Error{"file not found", "", nil} }
func MissingParam(param string) *Error { return &Error{"missing param", param, nil} }
func InvalidSignature() *Error         { return &Error{"invalid signature", "", nil} }
func InvalidAuthentication() *Error    { return &Error{"invalid authentication", "", nil} }

func (e *Error) Error() string {
	return e.Message
}
