package keygenext

import "github.com/keygen-sh/jsonapi-go"

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
