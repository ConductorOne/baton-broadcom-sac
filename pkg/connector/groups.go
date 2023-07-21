package connector

import (
	"context"
	"fmt"

	"github.com/ConductorOne/baton-broadcom-sac/pkg/sac"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	grant "github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type groupBuilder struct {
	resourceType *v2.ResourceType
	client       *sac.Client
}

const memberEntitlement = "member"

func (g *groupBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return g.resourceType
}

func groupResource(group *sac.Group, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"group_id":    group.ID,
		"group_name":  group.Name,
		"provider_id": group.IdentityProviderID,
	}

	groupTraitOptions := []rs.GroupTraitOption{rs.WithGroupProfile(profile)}

	ret, err := rs.NewGroupResource(
		group.Name,
		groupResourceType,
		group.ID,
		groupTraitOptions,
		rs.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (g *groupBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if parentResourceID == nil {
		return nil, "", nil, nil
	}

	groups, err := g.client.ListAllGroups(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to list groups: %w", err)
	}

	var rv []*v2.Resource
	for _, group := range groups {
		groupCopy := group
		gr, err := groupResource(&groupCopy, parentResourceID)
		if err != nil {
			return nil, "", nil, err
		}
		rv = append(rv, gr)
	}
	return rv, "", nil, nil
}

func (g *groupBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement

	assigmentOptions := []ent.EntitlementOption{
		ent.WithGrantableTo(userResourceType),
		ent.WithDescription(fmt.Sprintf("Member of %s Group", resource.DisplayName)),
		ent.WithDisplayName(fmt.Sprintf("%s Group %s", resource.DisplayName, memberEntitlement)),
	}

	en := ent.NewAssignmentEntitlement(resource, memberEntitlement, assigmentOptions...)
	rv = append(rv, en)

	return rv, "", nil, nil
}

func (g *groupBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var token string

	bag, err := parsePageToken(pToken.Token, resource.Id)
	if err != nil {
		return nil, "", nil, err
	}

	groupTrait, err := rs.GetGroupTrait(resource)
	if err != nil {
		return nil, "", nil, err
	}

	identityProviderId, ok := rs.GetProfileStringValue(groupTrait.Profile, "provider_id")
	if !ok {
		return nil, "", nil, fmt.Errorf("error fetching provider_id from group profile")
	}

	members, paginationData, err := g.client.ListGroupMembers(ctx, identityProviderId, resource.Id.Resource, bag.PageToken())
	if err != nil {
		return nil, "", nil, err
	}

	if !paginationData.Last {
		token, err = bag.NextToken(fmt.Sprintf("%v", paginationData.NextPage))
		if err != nil {
			return nil, "", nil, err
		}
	}

	var rv []*v2.Grant
	for _, member := range members {
		memberCopy := member
		ur, err := userResource(&memberCopy, resource.Id)
		if err != nil {
			return nil, "", nil, err
		}

		grant := grant.NewGrant(resource, memberEntitlement, ur.Id)
		rv = append(rv, grant)
	}

	return rv, token, nil, nil
}

func newGroupBuilder(client *sac.Client) *groupBuilder {
	return &groupBuilder{
		resourceType: groupResourceType,
		client:       client,
	}
}
