package connector

import (
	"context"
	"strings"

	sac "github.com/conductorone/baton-broadcom-sac/pkg/sac"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	resourceType *v2.ResourceType
	client       *sac.Client
}

func (u *userBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return u.resourceType
}

func baseUserResource(entity *sac.DirectoryEntity, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	names := strings.SplitN(entity.DisplayName, " ", 2)
	var firstName, lastName string
	switch len(names) {
	case 1:
		firstName = names[0]
	case 2:
		firstName = names[0]
		lastName = names[1]
	}

	profile := map[string]interface{}{
		"first_name":        valOrFallback(firstName, entity.DisplayName),
		"last_name":         valOrFallback(lastName, ""),
		"user_id":           entity.ID,
		"identity_provider": entity.IdentityProviderType,
	}

	userTraitOptions := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
	}

	ret, err := rs.NewUserResource(
		entity.DisplayName,
		userResourceType,
		entity.ID,
		userTraitOptions,
		rs.WithParentResourceID(parentResourceID),
	)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func userResource(user *sac.User, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"first_name":        valOrFallback(user.FirstName, user.Username),
		"last_name":         valOrFallback(user.LastName, ""),
		"login":             user.Email,
		"user_id":           user.ID,
		"identity_provider": user.RepositoryType,
	}

	var userStatus v2.UserTrait_Status_Status
	switch {
	case user.Blocked:
		userStatus = v2.UserTrait_Status_STATUS_DISABLED
	case user.IsDeleted:
		userStatus = v2.UserTrait_Status_STATUS_DELETED
	default:
		userStatus = v2.UserTrait_Status_STATUS_ENABLED
	}

	userTraitOptions := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithEmail(user.Email, true),
		rs.WithStatus(userStatus),
	}

	ret, err := rs.NewUserResource(
		user.Username,
		userResourceType,
		user.ID,
		userTraitOptions,
		rs.WithParentResourceID(parentResourceID),
	)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (u *userBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if parentResourceID == nil {
		return nil, "", nil, nil
	}

	var rv []*v2.Resource
	users, err := u.client.ListAllUsers(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	for _, user := range users {
		userCopy := user
		ur, err := userResource(&userCopy, parentResourceID)
		if err != nil {
			return nil, "", nil, err
		}
		rv = append(rv, ur)
	}

	return rv, "", nil, nil
}

func (o *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *userBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newUserBuilder(client *sac.Client) *userBuilder {
	return &userBuilder{
		resourceType: userResourceType,
		client:       client,
	}
}
