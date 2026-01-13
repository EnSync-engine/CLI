package domain

import "time"

const (
	permissionWildcard = "*"
)

type AccessKey struct {
	ID             string          `json:"id,omitempty"`
	AccessKey      string          `json:"accessKey"`
	Name           string          `json:"name,omitempty"`
	Type           string          `json:"type,omitempty"`
	CreatedAt      time.Time       `json:"createdAt,omitempty"`
	Permissions    *Permissions    `json:"permissions,omitempty"`
	ServiceKeyID   string          `json:"service_key_id,omitempty"`
	ServiceKeyPair *ServiceKeyPair `json:"service_key_pair,omitempty"`
}

type ServiceKeyPair struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key,omitempty"`
}

type Permissions struct {
	Send      []string               `json:"send,omitempty"`
	Receive   []string               `json:"receive,omitempty"`
	Resources map[string]interface{} `json:"resources,omitempty"`
}

type AccessKeyPermissions struct {
	ID             string          `json:"id,omitempty"`
	Key            string          `json:"key"`
	Name           string          `json:"name,omitempty"`
	Type           string          `json:"type,omitempty"`
	CreatedAt      time.Time       `json:"createdAt,omitempty"`
	Permissions    *Permissions    `json:"permissions"`
	ServiceKeyID   string          `json:"service_key_id,omitempty"`
	ServiceKeyPair *ServiceKeyPair `json:"service_key_pair,omitempty"`
}

type AccessKeyList struct {
	ResultsLength int                     `json:"resultsLength"`
	Results       []*AccessKeyPermissions `json:"results"`
}

type CreateAccessKeyRequest struct {
	Type        string       `json:"type"`
	Name        string       `json:"name"`
	Permissions *Permissions `json:"permissions"`
}

type UpdateServiceKeyPairRequest struct {
	AccessKey string `json:"access_key"`
}

func (p *Permissions) HasSendPermission(channel string) bool {
	return containsOrWildcard(p.Send, channel)
}

func (p *Permissions) HasReceivePermission(channel string) bool {
	return containsOrWildcard(p.Receive, channel)
}

func containsOrWildcard(slice []string, value string) bool {
	for _, v := range slice {
		if v == permissionWildcard || v == value {
			return true
		}
	}
	return false
}
