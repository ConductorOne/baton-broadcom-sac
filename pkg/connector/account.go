package connector

import (
	"context"
	"fmt"

	sac "github.com/conductorone/baton-broadcom-sac/pkg/sac"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	grant "github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

const (
	viewer = "viewer"
	admin  = "admin"
)

var roles = []string{
	viewer,
	admin,
}

type accountBuilder struct {
	resourceType *v2.ResourceType
	client       *sac.Client
	tenant       string
}

func (a *accountBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return a.resourceType
}

func accountResource(tenant string, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	ret, err := rs.NewResource(
		tenant,
		accountResourceType,
		tenant,
		rs.WithAnnotation(
			&v2.ChildResourceType{ResourceTypeId: userResourceType.Id},
			&v2.ChildResourceType{ResourceTypeId: groupResourceType.Id},
			&v2.ChildResourceType{ResourceTypeId: policyResourceType.Id},
		),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (a *accountBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var rv []*v2.Resource
	ur, err := accountResource(a.tenant, parentResourceID)
	if err != nil {
		return nil, "", nil, err
	}

	rv = append(rv, ur)

	return rv, "", nil, nil
}

func (a *accountBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement

	for _, role := range roles {
		permissionOptions := []ent.EntitlementOption{
			ent.WithGrantableTo(userResourceType),
			ent.WithDescription(fmt.Sprintf("%s user role", role)),
			ent.WithDisplayName(fmt.Sprintf("%s role %s", resource.DisplayName, role)),
		}
		en := ent.NewPermissionEntitlement(resource, role, permissionOptions...)
		rv = append(rv, en)
	}

	return rv, "", nil, nil
}

func (a *accountBuilder) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var rv []*v2.Grant
	users, err := a.client.ListAllUsers(ctx)
	if err != nil {
		return rv, "", nil, err
	}

	for _, user := range users {
		userCopy := user
		ur, err := userResource(&userCopy, resource.Id)
		if err != nil {
			return nil, "", nil, err
		}

		entitlement := viewer
		if user.IsAdmin {
			entitlement = admin
		}

		roleGrant := grant.NewGrant(resource, entitlement, ur.Id)
		rv = append(rv, roleGrant)
	}

	return rv, "", nil, nil
}

func newAccountBuilder(client *sac.Client) *accountBuilder {
	return &accountBuilder{
		resourceType: accountResourceType,
		client:       client,
	}
}
