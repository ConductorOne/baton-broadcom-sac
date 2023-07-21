# baton-broadcom-sac
`baton-broadcom-sac` is a connector for Broadcom SAC built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with the Broadcom SAC API to sync data about users, groups and roles.

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Getting Started

## Prerequisites

// TODO - add prerequisites 

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-broadcom-sac
baton-broadcom-sac
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_USERNAME=username BATON_PASSWORD=password BATON_TENANT=yourTenant ghcr.io/conductorone/baton-broadcom-sac:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-broadcom-sac/cmd/baton-broadcom-sac@main

BATON_USERNAME=username BATON_PASSWORD=password BATON_TENANT=yourTenant
baton resources
```

# Data Model

`baton-broadcom-sac` pulls down information about the following Broadcom SAC resources:
- Users
- Groups

# Contributing, Support, and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets. We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone. If you have questions, problems, or ideas: Please open a Github Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-broadcom-sac` Command Line Usage

```
baton-broadcom-sac

Usage:
  baton-broadcom-sac [flags]
  baton-broadcom-sac [command]

Available Commands:
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --client-id string              The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string          The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                   The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
      --grant-entitlement string      The entitlement to grant to the supplied principal ($BATON_GRANT_ENTITLEMENT)
      --grant-principal string        The resource to grant the entitlement to ($BATON_GRANT_PRINCIPAL)
      --grant-principal-type string   The resource type of the principal to grant the entitlement to ($BATON_GRANT_PRINCIPAL_TYPE)
  -h, --help                          help for baton-example
      --log-format string             The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string              The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --password string               Password for your Broadcom SAC instance. ($BATON_PASSWORD)
      --revoke-grant string           The grant to revoke ($BATON_REVOKE_GRANT)
      --tenant string                 Name of your Broadcom SAC tenant. ($BATON_TENANT)
      --username string               Username for your Broadcom SAC instance. ($BATON_USERNAME)
  -v, --version                       version for baton-example

Use "baton-broadcom-sac [command] --help" for more information about a command.

```
