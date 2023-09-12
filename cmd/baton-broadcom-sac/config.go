package main

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/spf13/cobra"
)

// config defines the external configuration required for the connector to run.
type config struct {
	cli.BaseConfig `mapstructure:",squash"` // Puts the base config options in the same place as the connector options

	SacClientID     string `mapstructure:"sac-client-id"`
	SacClientSecret string `mapstructure:"sac-client-secret"`
	Tenant          string `mapstructure:"tenant"`
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	if cfg.SacClientID == "" {
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

func cmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("sac-client-id", "", "Client ID for your Broadcom SAC instance. ($BATON_SAC_CLIENT_ID)")
	cmd.PersistentFlags().String("sac-client-secret", "", "Client Secret for your Broadcom SAC instance. ($BATON_SAC_CLIENT_SECRET)")
	cmd.PersistentFlags().String("tenant", "", "Name of your Broadcom SAC tenant. ($BATON_TENANT)")
}
