package keygenext

var (
	Account string
	Product string
	Token   string
)

type APIError struct {
	Title  string
	Detail string
	Code   string
	Err    error
}

func (e *APIError) Error() string {
	return e.Title + ": " + e.Detail
}

func (e *APIError) Unwrap() error {
	return e.Err
}
