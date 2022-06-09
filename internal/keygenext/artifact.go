package keygenext

import (
	"errors"
	"io"
	"net/http"

	"github.com/keygen-sh/jsonapi-go"
	"github.com/keygen-sh/keygen-go/v2"
)

// Artifact represents a Keygen artifact object.
type Artifact struct {
	ID        string `json:"-"`
	Type      string `json:"-"`
	Filename  string `json:"filename,omitempty"`
	Filetype  string `json:"filetype,omitempty"`
	Filesize  int64  `json:"filesize,omitempty"`
	Platform  string `json:"platform,omitempty"`
	Arch      string `json:"arch,omitempty"`
	Signature string `json:"signature,omitempty"`
	Checksum  string `json:"checksum,omitempty"`
	ReleaseID string `json:"-"`

	url string `json:"-"`
}

func (a *Artifact) SetID(id string) error {
	a.ID = id
	return nil
}

func (a *Artifact) SetType(t string) error {
	a.Type = t
	return nil
}

func (a *Artifact) SetData(to func(target interface{}) error) error {
	return to(a)
}

func (a Artifact) GetID() string {
	return a.ID
}

func (r Artifact) GetType() string {
	return "artifacts"
}

func (a Artifact) GetData() interface{} {
	return a
}

func (a Artifact) GetRelationships() map[string]interface{} {
	relationships := make(map[string]interface{})

	relationships["release"] = jsonapi.ResourceObjectIdentifier{
		Type: "releases",
		ID:   a.ReleaseID,
	}

	return relationships
}

func (a *Artifact) Create() error {
	client := keygen.NewClientWithOptions(
		&keygen.ClientOptions{Account: Account, Token: Token, PublicKey: PublicKey, UserAgent: UserAgent},
	)

	res, err := client.Post("artifacts", a, a)
	if err != nil {
		if len(res.Document.Errors) > 0 {
			e := res.Document.Errors[0]

			return &Error{Title: e.Title, Detail: e.Detail, Source: e.Source.Pointer, Code: e.Code, Err: err}
		}

		return err
	}

	a.url = res.Headers.Get("Location")

	return nil
}

func (a *Artifact) Upload(reader io.Reader) error {
	client := &http.Client{}

	req, err := http.NewRequest("PUT", a.url, reader)
	if err != nil {
		return err
	}

	// This must be set otherwise the Go http package sends a Transfer-Encoding
	// header, which S3 does not support.
	req.ContentLength = a.Filesize

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("failed to upload to storage provider")
	}

	return nil
}
