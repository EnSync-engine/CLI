package api

import (
	"context"

	"github.com/EnSync-engine/CLI/app/domain"
)

type APIClient interface {
	EventService
	AccessKeyService
	WorkspaceService
}

type EventService interface {
	ListEvents(ctx context.Context, params *ListParams) (*domain.EventList, error)
	GetEventByName(ctx context.Context, name string) (*domain.Event, error)
	CreateEvent(ctx context.Context, event *domain.Event) error
	UpdateEvent(ctx context.Context, event *domain.Event) error
}

type AccessKeyService interface {
	ListAccessKeys(ctx context.Context, params *ListParams) (*domain.AccessKeyList, error)
	GetAccessKeyByID(ctx context.Context, id string) (*domain.AccessKeyPermissions, error)
	CreateAccessKey(ctx context.Context, req *domain.CreateAccessKeyRequest) (*domain.AccessKey, error)
	DeleteAccessKey(ctx context.Context, id string) error
	GetAccessKeyPermissions(ctx context.Context, key string) (*domain.AccessKeyPermissions, error)
	SetAccessKeyPermissions(ctx context.Context, key string, permissions *domain.Permissions) error
	UpdateServiceKeyPair(ctx context.Context, accessKey string) (*domain.ServiceKeyPair, error)
}

type WorkspaceService interface {
	ListWorkspaces(ctx context.Context, params *ListParams) (*domain.WorkspaceList, error)
	CreateWorkspace(ctx context.Context, name string) error
}
