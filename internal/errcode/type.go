package errcode

// Error is a stable API error definition returned by services.
type Error struct {
	code   int
	reason string
}

func newError(code int, reason string) Error {
	return Error{
		code:   code,
		reason: reason,
	}
}

// Code returns the stable numeric API code.
func (e Error) Code() int {
	if e.code <= 0 {
		return 0
	}
	return e.code
}

// Reason returns the stable machine-readable error reason.
func (e Error) Reason() string {
	if e.reason == "" {
		return "UNKNOWN_ERROR"
	}
	return e.reason
}

// Error implements the standard error interface.
func (e Error) Error() string {
	return e.Reason()
}
