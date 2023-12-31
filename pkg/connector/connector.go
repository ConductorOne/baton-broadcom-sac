package connector

import (
	"context"
	"fmt"
	"io"

	sac "github.com/conductorone/baton-broadcom-sac/pkg/sac"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type Connector struct {
	client       *sac.Client
	clientID     string
	clientSecret string
	tenant       string
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from the upstream service.
func (c *Connector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newAccountBuilder(c.client),
		newUserBuilder(c.client),
		newGroupBuilder(c.client),
		newPolicyBuilder(c.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (c *Connector) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (c *Connector) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Broadcom SAC",
		Description: "Connector syncing users and groups from Broadcom SAC.",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (c *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	token, err := sac.CreateBearerToken(ctx, c.clientID, c.clientSecret, c.tenant)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	if token == "" {
		return nil, fmt.Errorf("missing access token")
	}

	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, clientID, clientSecret, tenant string) (*Connector, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	token, err := sac.CreateBearerToken(ctx, clientID, clientSecret, tenant)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	return &Connector{
		client:       sac.NewClient(httpClient, tenant, token),
		clientID:     clientID,
		clientSecret: clientSecret,
		tenant:       tenant,
	}, nil
}
