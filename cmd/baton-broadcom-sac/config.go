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

	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Tenant   string `mapstructure:"tenant"`
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	if cfg.Username == "" {
		return fmt.Errorf("username is missing")
	}

	if cfg.Password == "" {
		return fmt.Errorf("password is missing")
	}

	if cfg.Tenant == "" {
		return fmt.Errorf("tenant name is missing")
	}
	return nil
}

func cmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("username", "", "Username for your Broadcom SAC instance. ($BATON_USERNAME)")
	cmd.PersistentFlags().String("password", "", "Password for your Broadcom SAC instance. ($BATON_PASSWORD)")
	cmd.PersistentFlags().String("tenant", "", "Name of your Broadcom SAC tenant. ($BATON_TENANT)")
}
