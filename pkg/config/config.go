package config

import (
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	SacClientId = field.StringField(
		"sac-client-id",
		field.WithDescription("Client ID for your Broadcom SAC instance"),
		field.WithRequired(true),
		field.WithDisplayName("SAC Client ID"),
	)

	SacClientSecret = field.StringField(
		"sac-client-secret",
		field.WithDescription("Client Secret for your Broadcom SAC instance"),
		field.WithRequired(true),
		field.WithIsSecret(true),
		field.WithDisplayName("SAC Client Secret"),
	)

	Tenant = field.StringField(
		"tenant",
		field.WithDescription("Name of your Broadcom SAC tenant"),
		field.WithRequired(true),
		field.WithDisplayName("Tenant"),
	)

	// FieldRelationships defines relationships between the fields listed in
	// Config that can be automatically validated.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run ./gen
var Config = field.NewConfiguration([]field.SchemaField{
	SacClientId,
	SacClientSecret,
	Tenant,
})

// ValidateConfig is run after the configuration is loaded, and should return an
// error if it isn't valid. Implementing this function is optional, it only
// needs to perform extra validations that cannot be encoded with configuration
// parameters.
func ValidateConfig(cfg *BroadcomSac) error {
	if cfg.SacClientId == "" {
		return fmt.Errorf("client ID is missing")
	}

	if cfg.SacClientSecret == "" {
		return fmt.Errorf("client secret is missing")
	}

	if cfg.Tenant == "" {
		return fmt.Errorf("tenant name is missing")
	}

	return nil
}
