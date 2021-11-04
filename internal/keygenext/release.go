package keygenext

import (
	"os"

	"github.com/keygen-sh/jsonapi-go"
	"github.com/keygen-sh/keygen-go"
)

type Release struct {
	ID          string                 `json:"-"`
	Type        string                 `json:"-"`
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	Version     string                 `json:"version"`
	Filename    string                 `json:"filename"`
	Filetype    string                 `json:"filetype"`
	Filesize    int64                  `json:"filesize"`
	Platform    string                 `json:"platform"`
	Channel     string                 `json:"channel"`
	Signature   string                 `json:"signature"`
	Checksum    string                 `json:"checksum"`
	Metadata    map[string]interface{} `json:"metadata"`
	ProductID   string                 `json:"-"`
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

	relationships["constraints"] = r.Constraints
	relationships["product"] = jsonapi.ResourceObjectIdentifier{
		Type: "products",
		ID:   r.ProductID,
	}

	return relationships
}

func (r *Release) Upsert() error {
	client := &keygen.Client{Account: Account, Token: Token}

	res, err := client.Put("releases", r, r)
	if err != nil {
		if len(res.Document.Errors) > 0 {
			e := res.Document.Errors[0]

			return &APIError{Title: e.Title, Detail: e.Detail, Source: e.Source.Pointer, Code: e.Code, Err: err}
		}

		return err
	}

	return nil
}

func (r *Release) Upload(file *os.File) error {
	client := &keygen.Client{Account: Account, Token: Token}
	artifact := &Artifact{}

	res, err := client.Put("releases/"+r.ID+"/artifact", nil, artifact)
	if err != nil {
		if len(res.Document.Errors) > 0 {
			e := res.Document.Errors[0]

			return &APIError{Title: e.Title, Detail: e.Detail, Source: e.Source.Pointer, Code: e.Code, Err: err}
		}

		return err
	}

	artifact.ContentLength = r.Filesize
	artifact.Location = res.Headers.Get("Location")

	err = artifact.Upload(file)
	if err != nil {
		return err
	}

	return nil
}
