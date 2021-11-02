package keygenext

import "strings"

// FIXME(ezekg) Add overrides for UserAgent and PublicKey. Currently,
//              adding auto-upgrades to the CLI using our Go SDK will
//              cause issues since PublicKey is a global.
var (
	Account string
	Product string
	Token   string
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
