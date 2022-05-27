package keygenext

type Error struct {
	Title  string
	Detail string
	Code   string
	Source string
	Err    error
}

func (e *Error) Error() string {
	if e.Code != "" {
		return "[" + e.Code + "] " + e.Title + ": " + e.Detail
	} else {
		return e.Title + ": " + e.Detail
	}
}

func (e *Error) Unwrap() error {
	return e.Err
}
