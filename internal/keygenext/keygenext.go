package keygenext

import "strings"

var (
	Account   string
	Product   string
	Token     string
	PublicKey string
	UserAgent string
)

type APIError struct {
	Title  string
	Detail string
	Code   string
	Source string
	Err    error
}

func (e *APIError) Error() string {
	return strings.ToLower(e.Title) + ": " + e.Detail
}

func (e *APIError) Unwrap() error {
	return e.Err
}
