package keygenext

import (
	"github.com/keygen-sh/keygen-go/v2"
)

type Entitlement struct {
	keygen.Entitlement
}

func (e *Entitlement) Get() error {
	client := keygen.NewClientWithOptions(
		&keygen.ClientOptions{Account: Account, Environment: Environment, Token: Token, PublicKey: PublicKey, UserAgent: UserAgent},
	)

	res, err := client.Get("entitlements/"+e.ID, nil, e)
	if err != nil {
		if len(res.Document.Errors) > 0 {
			e := res.Document.Errors[0]

			return &Error{Title: e.Title, Detail: e.Detail, Source: e.Source.Pointer, Code: e.Code, Err: err}
		}

		return err
	}

	return nil
}
