// Package gcs provides a client for the Globus Connect Server Manager API.
package gcs

import "time"

// Endpoint represents a GCS endpoint configuration.
type Endpoint struct {
	ID                  string    `json:"id,omitempty"`
	DisplayName         string    `json:"display_name,omitempty"`
	Organization        string    `json:"organization,omitempty"`
	Department          string    `json:"department,omitempty"`
	Description         string    `json:"description,omitempty"`
	ContactEmail        string    `json:"contact_email,omitempty"`
	ContactInfo         string    `json:"contact_info,omitempty"`
	InfoLink            string    `json:"info_link,omitempty"`
	Public              bool      `json:"public,omitempty"`
	DefaultDirectory    string    `json:"default_directory,omitempty"`
	Keywords            []string  `json:"keywords,omitempty"`
	SubscriptionID      string    `json:"subscription_id,omitempty"`
	NetworkUse          string    `json:"network_use,omitempty"`
	MaxConcurrency      int       `json:"max_concurrency,omitempty"`
	PreferredConcurrency int      `json:"preferred_concurrency,omitempty"`
	DisableAnonymousWrites bool   `json:"disable_anonymous_writes,omitempty"`
	LastModified        time.Time `json:"last_modified,omitempty"`
}

// Info represents the GCS Manager service information.
type Info struct {
	APIVersion     string `json:"api_version"`
	EndpointID     string `json:"endpoint_id"`
	ManagerVersion string `json:"manager_version"`
}

// Collection represents a collection on a GCS endpoint.
type Collection struct {
	ID                  string            `json:"id,omitempty"`
	DisplayName         string            `json:"display_name,omitempty"`
	Description         string            `json:"description,omitempty"`
	CollectionType      string            `json:"collection_type,omitempty"`
	StorageGatewayID    string            `json:"storage_gateway_id,omitempty"`
	CollectionBaseFolder string           `json:"collection_base_path,omitempty"`
	Public              bool              `json:"public,omitempty"`
	DisableAnonymousWrites bool           `json:"disable_anonymous_writes,omitempty"`
	ContactEmail        string            `json:"contact_email,omitempty"`
	ContactInfo         string            `json:"contact_info,omitempty"`
	InfoLink            string            `json:"info_link,omitempty"`
	Keywords            []string          `json:"keywords,omitempty"`
	Organization        string            `json:"organization,omitempty"`
	Department          string            `json:"department,omitempty"`
	UserMessage         string            `json:"user_message,omitempty"`
	UserMessageLink     string            `json:"user_message_link,omitempty"`
	IdentityID          string            `json:"identity_id,omitempty"`
	Policies            *CollectionPolicies `json:"policies,omitempty"`
}

// CollectionPolicies represents access policies for a collection.
type CollectionPolicies struct {
	AuthenticationTimeoutMins int    `json:"authentication_timeout_mins,omitempty"`
	SharingRestrict           string `json:"sharing_restrict,omitempty"`
	SharingUsersAllow         []string `json:"sharing_users_allow,omitempty"`
	SharingUsersDeny          []string `json:"sharing_users_deny,omitempty"`
	SharingGroupsAllow        []string `json:"sharing_groups_allow,omitempty"`
	SharingGroupsDeny         []string `json:"sharing_groups_deny,omitempty"`
}

// StorageGateway represents a storage backend configuration.
type StorageGateway struct {
	ID                  string            `json:"id,omitempty"`
	DisplayName         string            `json:"display_name,omitempty"`
	ConnectorID         string            `json:"connector_id,omitempty"`
	ConnectorName       string            `json:"connector_name,omitempty"`
	Root                string            `json:"root,omitempty"`
	IdentityMappings    []IdentityMapping `json:"identity_mappings,omitempty"`
	AllowedDomains      []string          `json:"allowed_domains,omitempty"`
	HighAssurance       bool              `json:"high_assurance,omitempty"`
	RequireMFA          bool              `json:"require_mfa,omitempty"`
	RestrictPaths       *PathRestrictions `json:"restrict_paths,omitempty"`
	PosixStagingFolder  string            `json:"posix_staging_path,omitempty"`
	PosixUserIDMap      string            `json:"posix_user_id_map,omitempty"`
	PosixGroupIDMap     string            `json:"posix_group_id_map,omitempty"`
	Policies            *StorageGatewayPolicies `json:"policies,omitempty"`
}

// IdentityMapping represents an identity mapping for a storage gateway.
type IdentityMapping struct {
	DataAccessProtocol string `json:"data_access_protocol,omitempty"`
	IdentityID         string `json:"identity_id,omitempty"`
	LocalUsername      string `json:"local_username,omitempty"`
}

// PathRestrictions represents path access restrictions.
type PathRestrictions struct {
	ReadOnly  []string `json:"read_only,omitempty"`
	ReadWrite []string `json:"read_write,omitempty"`
	None      []string `json:"none,omitempty"`
}

// StorageGatewayPolicies represents policies for a storage gateway.
type StorageGatewayPolicies struct {
	DataType string `json:"DATA_TYPE,omitempty"`
}

// Role represents an access role assignment.
type Role struct {
	ID         string `json:"id,omitempty"`
	Collection string `json:"collection,omitempty"`
	Principal  string `json:"principal,omitempty"`
	Role       string `json:"role,omitempty"`
}

// Node represents a GCS node configuration.
type Node struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Incoming  bool   `json:"incoming,omitempty"`
	Outgoing  bool   `json:"outgoing,omitempty"`
}

// DomainConfig represents custom domain configuration.
type DomainConfig struct {
	Domain      string `json:"domain"`
	Certificate string `json:"certificate,omitempty"`
	PrivateKey  string `json:"private_key,omitempty"`
	Verified    bool   `json:"verified,omitempty"`
}
