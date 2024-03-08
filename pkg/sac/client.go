package sac

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type Client struct {
	httpClient *http.Client
	baseUrl    string
	token      string
}

func NewClient(httpClient *http.Client, tenant, token string) *Client {
	baseUrl := fmt.Sprintf("https://api.%s.luminatesec.com/v2", tenant)
	return &Client{
		httpClient: httpClient,
		baseUrl:    baseUrl,
		token:      token,
	}
}

type AuthResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type PaginationData struct {
	First            bool   `json:"first"`
	Last             bool   `json:"last"`
	Size             int    `json:"size"`
	TotalElements    int    `json:"totalElements"`
	PerPage          int    `json:"perPage"`
	NextPage         string `json:"nextPage"`
	TotalPages       int    `json:"totalPages"`
	Number           int    `json:"number"`
	NumberOfElements int    `json:"numberOfElements"`
}

const applicationJSONHeader = "application/json"

// returns query params with pagination options with PageOffset.
func paginationQueryOffset(nextPage string) url.Values {
	q := url.Values{}
	if nextPage != "" {
		q.Set("pageOffset", nextPage)
	}
	q.Set("perPage", "50")
	return q
}

// returns query params with pagination options with PageNumber.
func paginationQueryPages(pageNumber int) url.Values {
	q := url.Values{}
	q.Set("page", fmt.Sprintf("%v", pageNumber))
	q.Set("size", "50")
	return q
}

// CreateBearerToken creates a bearer token for the given username, password, and tenant.
func CreateBearerToken(ctx context.Context, username, password, tenant string) (string, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.%s.luminatesec.com/v1/oauth/token", tenant)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", applicationJSONHeader)
	req.Header.Add("Content-Type", applicationJSONHeader)
	req.SetBasicAuth(username, password)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var res AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if res.Error != "" {
		return "", fmt.Errorf("error creating bearer token. %s: %s", res.Error, res.ErrorDescription)
	}

	return res.AccessToken, nil
}

// ListIdentityProviderIDs returns a list of identity provider ids.
func (c *Client) ListIdentityProviderIDs(ctx context.Context) ([]string, error) {
	var providerIDs []string
	providersUrl := fmt.Sprintf("%s/identities/settings/identity-providers", c.baseUrl)

	q := url.Values{}
	q.Set("includeLocal", "true")

	var res []IdentityProvider
	if err := c.doRequest(ctx, providersUrl, &res, q); err != nil {
		return nil, err
	}

	for _, identityProvider := range res {
		providerIDs = append(providerIDs, identityProvider.ID)
	}

	return providerIDs, nil
}

// ListUserPerProvider returns a list of users for the given identity provider id.
func (c *Client) ListUsersPerProvider(ctx context.Context, identityProviderId string, nextPage string) ([]User, PaginationData, error) {
	url := fmt.Sprintf("%s/identities/%s/users", c.baseUrl, identityProviderId)
	var res struct {
		Content []User `json:"content"`
		PaginationData
	}

	q := paginationQueryOffset(nextPage)
	if err := c.doRequest(ctx, url, &res, q); err != nil {
		return nil, PaginationData{}, err
	}

	return res.Content, res.PaginationData, nil
}

// ListAllUsers returns a list of all users for all identity providers.
func (c *Client) ListAllUsers(ctx context.Context) ([]User, error) {
	var allUsers []User
	identityProviders, err := c.ListIdentityProviderIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching identity providers: %w", err)
	}
	for _, identityProvider := range identityProviders {
		var nextPage string
		for {
			users, paginationData, err := c.ListUsersPerProvider(ctx, identityProvider, nextPage)
			if err != nil {
				return nil, fmt.Errorf("error fetching users: %w", err)
			}

			allUsers = append(allUsers, users...)
			if paginationData.Last {
				break
			}
			nextPage = paginationData.NextPage
		}
	}

	return allUsers, nil
}

// ListGroups returns a list of groups for the given identity provider id.
func (c *Client) ListGroupsPerProvider(ctx context.Context, identityProviderId string, nextPage string) ([]Group, PaginationData, error) {
	url := fmt.Sprintf("%s/identities/%s/groups", c.baseUrl, identityProviderId)
	var res struct {
		Content []Group `json:"content"`
		PaginationData
	}

	q := paginationQueryOffset(nextPage)

	if err := c.doRequest(ctx, url, &res, q); err != nil {
		return nil, PaginationData{}, err
	}

	return res.Content, res.PaginationData, nil
}

// ListAllGroups returns a list of all groups for all identity providers.
func (c *Client) ListAllGroups(ctx context.Context) ([]Group, error) {
	var allGroups []Group
	identityProviders, err := c.ListIdentityProviderIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching identity providers: %w", err)
	}
	for _, identityProvider := range identityProviders {
		var nextPage string
		for {
			groups, paginationData, err := c.ListGroupsPerProvider(ctx, identityProvider, nextPage)
			if err != nil {
				return nil, fmt.Errorf("error fetching groups: %w", err)
			}

			allGroups = append(allGroups, groups...)

			if paginationData.Last {
				break
			}

			nextPage = paginationData.NextPage
		}
	}

	return allGroups, nil
}

// ListGroupUsers returns a list of users for the given identity provider id and group id.
func (c *Client) ListGroupMembers(ctx context.Context, identityProviderId string, groupId string, nextPage string) ([]User, PaginationData, error) {
	url := fmt.Sprintf("%s/identities/%s/groups/%s/users", c.baseUrl, identityProviderId, groupId)
	var res struct {
		Content []User `json:"content"`
		PaginationData
	}

	q := paginationQueryOffset(nextPage)

	if err := c.doRequest(ctx, url, &res, q); err != nil {
		return nil, PaginationData{}, err
	}

	return res.Content, res.PaginationData, nil
}

// List Policies returns a list of policies.
func (c *Client) ListPolicies(ctx context.Context, pageNumber int) ([]Policy, PaginationData, error) {
	url := fmt.Sprintf("%s/policies", c.baseUrl)
	var res struct {
		Content []Policy `json:"content"`
		PaginationData
	}

	q := paginationQueryPages(pageNumber)

	if err := c.doRequest(ctx, url, &res, q); err != nil {
		return nil, PaginationData{}, err
	}

	return res.Content, res.PaginationData, nil
}

// ListAllPolicies returns a paginated list of all policies.
func (c *Client) ListAllPolicies(ctx context.Context) ([]Policy, error) {
	var nextPage int
	var allPolicies []Policy
	for {
		policies, paginationData, err := c.ListPolicies(ctx, nextPage)
		if err != nil {
			return nil, fmt.Errorf("error fetching policies: %w", err)
		}

		allPolicies = append(allPolicies, policies...)

		if paginationData.Last {
			break
		}

		nextPage = paginationData.Number + 1
	}

	return allPolicies, nil
}

// GetPolicy returns a policy by ID.
func (c *Client) GetPolicy(ctx context.Context, policyId string) (Policy, error) {
	url := fmt.Sprintf("%s/policies/%s", c.baseUrl, policyId)
	var res Policy

	if err := c.doRequest(ctx, url, &res, nil); err != nil {
		return Policy{}, err
	}

	return res, nil
}

func (c *Client) doRequest(ctx context.Context, url string, res interface{}, query url.Values) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	req.Header.Add("Accept", applicationJSONHeader)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	return nil
}
