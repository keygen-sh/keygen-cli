package keygenext

import (
	"fmt"

	"github.com/keygen-sh/jsonapi-go"
	"github.com/keygen-sh/keygen-go"
)

type Release struct {
	ID          string                 `json:"-"`
	Type        string                 `json:"-"`
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Filename    string                 `json:"filename"`
	Filetype    string                 `json:"filetype"`
	Filesize    int                    `json:"filesize"`
	Platform    string                 `json:"platform"`
	Channel     string                 `json:"channel"`
	Metadata    map[string]interface{} `json:"metadata"`
	ProductID   string                 `json:"-"`
	Constraints Constraints            `json:"-"`
	Location    string                 `json:"-"`
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
	// TODO(ezekg) Add custom user agent
	logger := keygen.Logger.(*keygen.LeveledLogger)
	logger.Level = keygen.LogLevelDebug

	client := &keygen.Client{Account: Account, Token: Token}

	fmt.Printf("%v\n", r)

	res, err := client.Put("releases", r, r)
	if err != nil {
		fmt.Printf("%v\n", res.Document.Errors[0])

		return err
	}

	return nil
}
