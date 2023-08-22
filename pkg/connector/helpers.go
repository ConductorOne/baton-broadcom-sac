package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
)

func parsePageToken(token string, resourceID *v2.ResourceId) (*pagination.Bag, error) {
	b := &pagination.Bag{}
	err := b.Unmarshal(token)
	if err != nil {
		return nil, err
	}

	if b.Current() == nil {
		b.Push(pagination.PageState{
			ResourceTypeID: resourceID.ResourceType,
			ResourceID:     resourceID.Resource,
		})
	}

	return b, nil
}

func valOrFallback(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}
