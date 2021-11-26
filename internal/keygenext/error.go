package keygenext

type APIError struct {
	Title  string
	Detail string
	Code   string
	Source string
	Err    error
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return e.Code + " - " + e.Title + ": " + e.Detail
	} else {
		return e.Title + ": " + e.Detail
	}
}

func (e *APIError) Unwrap() error {
	return e.Err
}
