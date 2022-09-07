package keygenext

import (
	"github.com/google/uuid"
	"github.com/keygen-sh/jsonapi-go"
	"github.com/keygen-sh/keygen-go/v2"
)

type Constraint struct {
	ID            string `json:"-"`
	Type          string `json:"-"`
	EntitlementID string `json:"-"`
}

func (c Constraint) GetID() string {
	return c.ID
}

func (c Constraint) GetType() string {
	return "constraints"
}

func (c Constraint) GetData() interface{} {
	return c
}

func (c Constraint) GetRelationships() map[string]interface{} {
	relationships := make(map[string]interface{})

	relationships["entitlement"] = jsonapi.ResourceObjectIdentifier{
		Type: "entitlements",
		ID:   c.EntitlementID,
	}

	return relationships
}

func (c Constraint) UseExperimentalEmbeddedRelationshipData() bool {
	return true
}

type Constraints []Constraint

func (c Constraints) GetData() interface{} {
	return c
}

func (c Constraints) From(entitlements []string) Constraints {
	for _, identifier := range entitlements {
		if _, err := uuid.Parse(identifier); err != nil {
			entitlement := &Entitlement{keygen.Entitlement{ID: identifier}}

			// identifier may be an ID or an entitlement code, so we're
			// retrieving the entitlement to get it's real ID.
			entitlement.Get()

			identifier = entitlement.ID
		}

		c = append(c, Constraint{EntitlementID: identifier})
	}

	return c
}
