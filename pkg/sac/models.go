package sac

type User struct {
	Username           string `json:"username"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Email              string `json:"email"`
	RepositoryType     string `json:"repository_type"`
	IsAdmin            bool   `json:"is_admin"`
	IsDeleted          bool   `json:"is_deleted"`
	Blocked            bool   `json:"blocked"`
	NotificationEmail  string `json:"notification_email"`
	ID                 string `json:"id"`
	IdentityProviderID string `json:"identity_provider_id"`
}

type Group struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	RepositoryType     string `json:"repository_type"`
	IdentityProviderID string `json:"identity_provider_id"`
}

type IdentityProvider struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Provider        string      `json:"provider"`
	IsAuthenticator bool        `json:"is_authenticator"`
	IsUserStore     bool        `json:"is_user_store"`
	AuthenticatorID interface{} `json:"authenticator_id"`
}

type Account struct {
	Name string
	ID   string
}

type Policy struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Type              string            `json:"type"`
	Enabled           bool              `json:"enabled"`
	CreatedAt         string            `json:"createdAt"`
	DirectoryEntities []DirectoryEntity `json:"directoryEntities"`
	PolicyAccess      string            `json:"PolicyAccess"`
}

// Either User or Group.
type DirectoryEntity struct {
	ID                   string `json:"id"`
	IdentifierInProvider string `json:"identifierInProvider"`
	IdentityProviderID   string `json:"identityProviderId"`
	IdentityProviderType string `json:"identityProviderType"`
	Type                 string `json:"type"`
	DisplayName          string `json:"displayName"`
}
