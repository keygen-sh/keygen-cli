package keygenext

import (
	"github.com/google/go-querystring/query"
	"github.com/keygen-sh/jsonapi-go"
	"github.com/keygen-sh/keygen-go/v2"
)

type Release struct {
	ID          string                 `json:"-"`
	Type        string                 `json:"-"`
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Version     string                 `json:"version,omitempty"`
	Tag         *string                `json:"tag"`
	Channel     string                 `json:"channel,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	ProductID   string                 `json:"-"`
	PackageID   *string                `json:"-"`
	Constraints Constraints            `json:"-"`
}

func (r *Release) SetID(id string) error {
	r.ID = id
	return nil
}

func (r *Release) SetType(t string) error {
	r.Type = t
	return nil
}

func (r *Release) SetData(to func(target interface{}) error) error {
	return to(r)
}

func (r Release) GetID() string {
	return r.ID
}

func (r Release) GetType() string {
	return "releases"
}

func (r Release) GetData() interface{} {
	return r
}

func (r Release) GetRelationships() map[string]interface{} {
	relationships := make(map[string]interface{})

	if len(r.Constraints) > 0 {
		relationships["constraints"] = r.Constraints
	}

	if r.ProductID != "" {
		relationships["product"] = jsonapi.ResourceObjectIdentifier{
			Type: "products",
			ID:   r.ProductID,
		}
	}

	if r.PackageID != nil {
		relationships["package"] = jsonapi.ResourceObjectIdentifier{
			Type: "packages",
			ID:   *r.PackageID,
		}
	}

	if len(relationships) == 0 {
		return nil
	}

	return relationships
}

func (r *Release) Create() error {
	client := keygen.NewClientWithOptions(
		&keygen.ClientOptions{Account: Account, Environment: Environment, Token: Token, PublicKey: PublicKey, UserAgent: UserAgent, APIURL: APIURL},
	)

	res, err := client.Post("releases", r, r)
	if err != nil {
		if res != nil && len(res.Document.Errors) > 0 {
			e := res.Document.Errors[0]

			return &Error{Title: e.Title, Detail: e.Detail, Source: e.Source.Pointer, Code: e.Code, Err: err}
		}

		return err
	}

	return nil
}

func (r *Release) Update() error {
	client := keygen.NewClientWithOptions(
		&keygen.ClientOptions{Account: Account, Environment: Environment, Token: Token, PublicKey: PublicKey, UserAgent: UserAgent, APIURL: APIURL},
	)

	res, err := client.Patch("releases/"+r.ID, r, r)
	if err != nil {
		if res != nil && len(res.Document.Errors) > 0 {
			e := res.Document.Errors[0]

			return &Error{Title: e.Title, Detail: e.Detail, Source: e.Source.Pointer, Code: e.Code, Err: err}
		}

		return err
	}

	return nil
}

func (r *Release) Get() error {
	client := keygen.NewClientWithOptions(
		&keygen.ClientOptions{Account: Account, Environment: Environment, Token: Token, PublicKey: PublicKey, UserAgent: UserAgent, APIURL: APIURL},
	)

	// TODO(ezekg) Add support for custom query params to SDK
	type querystring struct {
		Package string `url:"package,omitempty"`
		Product string `url:"product,omitempty"`
	}

	qs := querystring{Package: *r.PackageID, Product: Product}
	values, err := query.Values(qs)
	if err != nil {
		return err
	}

	url := "releases/" + r.ID
	if enc := values.Encode(); enc != "" {
		url += "?" + enc
	}

	res, err := client.Get(url, nil, r)
	if err != nil {
		if res != nil && len(res.Document.Errors) > 0 {
			e := res.Document.Errors[0]

			return &Error{Title: e.Title, Detail: e.Detail, Source: e.Source.Pointer, Code: e.Code, Err: err}
		}

		return err
	}

	return nil
}

func (r *Release) Publish() error {
	client := keygen.NewClientWithOptions(
		&keygen.ClientOptions{Account: Account, Environment: Environment, Token: Token, PublicKey: PublicKey, UserAgent: UserAgent, APIURL: APIURL},
	)

	res, err := client.Post("releases/"+r.ID+"/actions/publish", nil, r)
	if err != nil {
		if res != nil && len(res.Document.Errors) > 0 {
			e := res.Document.Errors[0]

			return &Error{Title: e.Title, Detail: e.Detail, Source: e.Source.Pointer, Code: e.Code, Err: err}
		}

		return err
	}

	return nil
}

func (r *Release) Yank() error {
	client := keygen.NewClientWithOptions(
		&keygen.ClientOptions{Account: Account, Environment: Environment, Token: Token, PublicKey: PublicKey, UserAgent: UserAgent, APIURL: APIURL},
	)

	res, err := client.Post("releases/"+r.ID+"/actions/yank", nil, r)
	if err != nil {
		if res != nil && len(res.Document.Errors) > 0 {
			e := res.Document.Errors[0]

			return &Error{Title: e.Title, Detail: e.Detail, Source: e.Source.Pointer, Code: e.Code, Err: err}
		}

		return err
	}

	return nil
}

func (r *Release) Delete() error {
	client := keygen.NewClientWithOptions(
		&keygen.ClientOptions{Account: Account, Environment: Environment, Token: Token, PublicKey: PublicKey, UserAgent: UserAgent, APIURL: APIURL},
	)

	res, err := client.Delete("releases/"+r.ID, nil, r)
	if err != nil {
		if res != nil && len(res.Document.Errors) > 0 {
			e := res.Document.Errors[0]

			return &Error{Title: e.Title, Detail: e.Detail, Source: e.Source.Pointer, Code: e.Code, Err: err}
		}

		return err
	}

	return nil
}
