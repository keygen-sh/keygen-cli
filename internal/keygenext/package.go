package keygenext

import (
	"net/url"

	"github.com/keygen-sh/jsonapi-go"
	"github.com/keygen-sh/keygen-go/v2"
)

type Package struct {
	ID        string                 `json:"-"`
	Type      string                 `json:"-"`
	Name      *string                `json:"name,omitempty"`
	Key       *string                `json:"key,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	ProductID string                 `json:"-"`
}

func (p *Package) SetID(id string) error {
	p.ID = id
	return nil
}

func (p *Package) SetType(t string) error {
	p.Type = t
	return nil
}

func (p *Package) SetData(to func(target interface{}) error) error {
	return to(p)
}

func (p Package) GetID() string {
	return p.ID
}

func (p Package) GetType() string {
	return "packages"
}

func (p Package) GetData() interface{} {
	return p
}

func (p Package) GetRelationships() map[string]interface{} {
	relationships := make(map[string]interface{})

	if p.ProductID != "" {
		relationships["product"] = jsonapi.ResourceObjectIdentifier{
			Type: "products",
			ID:   p.ProductID,
		}
	}

	if len(relationships) == 0 {
		return nil
	}

	return relationships
}

func (p *Package) Get() error {
	client := keygen.NewClientWithOptions(
		&keygen.ClientOptions{Account: Account, Environment: Environment, Token: Token, PublicKey: PublicKey, UserAgent: UserAgent, APIURL: APIURL},
	)

	res, err := client.Get("packages/"+url.PathEscape(p.ID), nil, p)
	if err != nil {
		if res != nil && len(res.Document.Errors) > 0 {
			e := res.Document.Errors[0]

			return &Error{Title: e.Title, Detail: e.Detail, Source: e.Source.Pointer, Code: e.Code, Err: err}
		}

		return err
	}

	return nil
}
