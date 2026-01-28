package main

import (
	cfg "github.com/conductorone/baton-broadcom-sac/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("broadcom-sac", cfg.Config)
}
