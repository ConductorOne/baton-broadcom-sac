package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-broadcom-sac/pkg/sac"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	grant "github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type policyBuilder struct {
	resourceType *v2.ResourceType
	client       *sac.Client
}

const (
	assignmentEntitlement = "assigned"
	user                  = "User"
	group                 = "Group"
)

func (p *policyBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return p.resourceType
}

func policyResource(policy *sac.Policy, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	ret, err := rs.NewResource(
		policy.Name,
		policyResourceType,
		policy.ID,
		rs.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (p *policyBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if parentResourceID == nil {
		return nil, "", nil, nil
	}

	policies, err := p.client.ListAllPolicies(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to list policies: %w", err)
	}

	var rv []*v2.Resource
	for _, policy := range policies {
		policyCopy := policy
		gr, err := policyResource(&policyCopy, parentResourceID)
		if err != nil {
			return nil, "", nil, err
		}
		rv = append(rv, gr)
	}
	return rv, "", nil, nil
}

func (p *policyBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement

	assigmentOptions := []ent.EntitlementOption{
		ent.WithGrantableTo(userResourceType, groupResourceType),
		ent.WithDescription(fmt.Sprintf("Assigned to %s policy", resource.DisplayName)),
		ent.WithDisplayName(fmt.Sprintf("%s to %s policy", assignmentEntitlement, resource.DisplayName)),
	}

	en := ent.NewAssignmentEntitlement(resource, assignmentEntitlement, assigmentOptions...)
	rv = append(rv, en)

	return rv, "", nil, nil
}

func (p *policyBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	policy, err := p.client.GetPolicy(ctx, resource.Id.Resource)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to get policy: %w", err)
	}

	var rv []*v2.Grant
	for _, entity := range policy.DirectoryEntities {
		entityCopy := entity

		switch entityCopy.Type {
		case user:
			ur, err := baseUserResource(&entityCopy, resource.Id)
			if err != nil {
				return nil, "", nil, err
			}
			grant := grant.NewGrant(resource, assignmentEntitlement, ur.Id)
			rv = append(rv, grant)

		case group:
			gr, err := baseGroupResource(&entityCopy, resource.Id)
			if err != nil {
				return nil, "", nil, err
			}
			grant := grant.NewGrant(resource, assignmentEntitlement, gr.Id)
			rv = append(rv, grant)

		default:
			continue
		}
	}

	return rv, "", nil, nil
}

func newPolicyBuilder(client *sac.Client) *policyBuilder {
	return &policyBuilder{
		resourceType: policyResourceType,
		client:       client,
	}
}
