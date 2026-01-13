package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EnSync-engine/CLI/app/api"
	"github.com/EnSync-engine/CLI/app/domain"
)

const (
	testAccessKey = "test-access-key"
	headerAccess  = "X-ACCESS-KEY"
)

var (
	eventNamePattern     = regexp.MustCompile(`^/event/[\w-]+$`)
	accessKeyIDPattern   = regexp.MustCompile(`^/access-key/[\w-]+$`)
	accessKeyPermPattern = regexp.MustCompile(`^/access-key/[\w-]+/permissions$`)
)

func TestClient(t *testing.T) {
	server := newMockServer(t)
	defer server.Close()

	client := api.NewClient(server.URL)
	client.SetAccessKey(testAccessKey)

	ctx := context.Background()

	t.Run("Events", func(t *testing.T) {
		testEventOperations(t, ctx, client)
	})

	t.Run("AccessKeys", func(t *testing.T) {
		testAccessKeyOperations(t, ctx, client)
	})

	t.Run("Workspaces", func(t *testing.T) {
		testWorkspaceOperations(t, ctx, client)
	})
}

func testEventOperations(t *testing.T, ctx context.Context, client *api.Client) {
	t.Run("List", func(t *testing.T) {
		params := api.DefaultListParams()

		events, err := client.ListEvents(ctx, params)

		require.NoError(t, err)
		require.NotNil(t, events)
		assert.Equal(t, 2, events.ResultsLength)
		assert.Len(t, events.Results, 2)
	})

	t.Run("GetByName", func(t *testing.T) {
		event, err := client.GetEventByName(ctx, "test-event")

		require.NoError(t, err)
		require.NotNil(t, event)
		assert.Equal(t, "test-event", event.Name)
	})

	t.Run("Create", func(t *testing.T) {
		event := &domain.Event{
			Name:    "new-event",
			Payload: map[string]any{"key": "value"},
		}

		err := client.CreateEvent(ctx, event)

		require.NoError(t, err)
	})

	t.Run("Update", func(t *testing.T) {
		event := &domain.Event{
			ID:      "event-123",
			Name:    "updated-event",
			Payload: map[string]any{"key": "new-value"},
		}

		err := client.UpdateEvent(ctx, event)

		require.NoError(t, err)
	})
}

func testAccessKeyOperations(t *testing.T, ctx context.Context, client *api.Client) {
	t.Run("List", func(t *testing.T) {
		params := api.DefaultListParams()

		keys, err := client.ListAccessKeys(ctx, params)

		require.NoError(t, err)
		require.NotNil(t, keys)
		assert.Equal(t, 2, keys.ResultsLength)
	})

	t.Run("GetByID", func(t *testing.T) {
		key, err := client.GetAccessKeyByID(ctx, "test-id-123")

		require.NoError(t, err)
		require.NotNil(t, key)
		assert.Equal(t, "test-id-123", key.ID)
	})

	t.Run("Create", func(t *testing.T) {
		req := &domain.CreateAccessKeyRequest{
			Type: "SERVICE",
			Name: "test",
			Permissions: &domain.Permissions{
				Send:    []string{"event1", "event2"},
				Receive: []string{"event3"},
			},
		}

		created, err := client.CreateAccessKey(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.AccessKey)
	})

	t.Run("Delete", func(t *testing.T) {
		err := client.DeleteAccessKey(ctx, "delete-id-123")

		require.NoError(t, err)
	})

	t.Run("GetPermissions", func(t *testing.T) {
		perms, err := client.GetAccessKeyPermissions(ctx, "test-key")

		require.NoError(t, err)
		require.NotNil(t, perms)
		assert.NotNil(t, perms.Permissions)
	})

	t.Run("SetPermissions", func(t *testing.T) {
		permissions := &domain.Permissions{
			Send:    []string{"event1", "event2"},
			Receive: []string{"event3"},
		}

		err := client.SetAccessKeyPermissions(ctx, "test-key", permissions)

		require.NoError(t, err)
	})

	t.Run("UpdateServiceKeyPair", func(t *testing.T) {
		keyPair, err := client.UpdateServiceKeyPair(ctx, "test-key")

		require.NoError(t, err)
		require.NotNil(t, keyPair)
		assert.NotEmpty(t, keyPair.PublicKey)
		assert.NotEmpty(t, keyPair.PrivateKey)
	})
}

func testWorkspaceOperations(t *testing.T, ctx context.Context, client *api.Client) {
	t.Run("List", func(t *testing.T) {
		params := api.DefaultListParams()

		workspaces, err := client.ListWorkspaces(ctx, params)

		require.NoError(t, err)
		require.NotNil(t, workspaces)
		assert.Equal(t, 1, workspaces.ResultsLength)
	})

	t.Run("Create", func(t *testing.T) {
		err := client.CreateWorkspace(ctx, "new-workspace")

		require.NoError(t, err)
	})
}

func newMockServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		assert.Equal(t, testAccessKey, r.Header.Get(headerAccess))

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/event":
			handleListEvents(w)
		case r.Method == http.MethodPost && r.URL.Path == "/event":
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodGet && eventNamePattern.MatchString(r.URL.Path):
			handleGetEvent(w)
		case r.Method == http.MethodPut && eventNamePattern.MatchString(r.URL.Path):
			w.WriteHeader(http.StatusOK)

		case r.Method == http.MethodGet && r.URL.Path == "/access-key":
			handleListAccessKeys(w)
		case r.Method == http.MethodPost && r.URL.Path == "/access-key":
			handleCreateAccessKey(w)
		case r.Method == http.MethodGet && accessKeyIDPattern.MatchString(r.URL.Path) && !accessKeyPermPattern.MatchString(r.URL.Path):
			handleGetAccessKeyByID(w)
		case r.Method == http.MethodDelete && accessKeyIDPattern.MatchString(r.URL.Path):
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && accessKeyPermPattern.MatchString(r.URL.Path):
			handleGetAccessKeyPermissions(w)
		case r.Method == http.MethodPost && accessKeyPermPattern.MatchString(r.URL.Path):
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodPut && r.URL.Path == "/access/service-key-pair":
			handleUpdateServiceKeyPair(w)

		case r.Method == http.MethodGet && r.URL.Path == "/workspace":
			handleListWorkspaces(w)
		case r.Method == http.MethodPost && r.URL.Path == "/workspace":
			w.WriteHeader(http.StatusOK)
			writeJSON(w, map[string]any{})

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func handleListEvents(w http.ResponseWriter) {
	writeJSON(w, domain.EventList{
		ResultsLength: 2,
		Results: []*domain.Event{
			{ID: "event-1", Name: "event1", Payload: map[string]any{"key": "value1"}},
			{ID: "event-2", Name: "event2", Payload: map[string]any{"key": "value2"}},
		},
	})
}

func handleGetEvent(w http.ResponseWriter) {
	writeJSON(w, domain.Event{
		ID:      "event-1",
		Name:    "test-event",
		Payload: map[string]any{"key": "value"},
	})
}

func handleListAccessKeys(w http.ResponseWriter) {
	writeJSON(w, domain.AccessKeyList{
		ResultsLength: 2,
		Results: []*domain.AccessKeyPermissions{
			{ID: "id1", Key: "key1", Permissions: &domain.Permissions{Send: []string{"event1"}, Receive: []string{"event2"}}},
			{ID: "id2", Key: "key2", Permissions: &domain.Permissions{Send: []string{"event3"}, Receive: []string{"event4"}}},
		},
	})
}

func handleCreateAccessKey(w http.ResponseWriter) {
	writeJSON(w, domain.AccessKey{
		ID:        "new-id-123",
		AccessKey: "new-access-key-123",
		Name:      "test",
		Type:      "SERVICE",
	})
}

func handleGetAccessKeyByID(w http.ResponseWriter) {
	writeJSON(w, domain.AccessKeyPermissions{
		ID:          "test-id-123",
		Key:         "test-key",
		Name:        "Test Key",
		Type:        "SERVICE",
		Permissions: &domain.Permissions{Send: []string{"event1"}, Receive: []string{"event2"}},
	})
}

func handleGetAccessKeyPermissions(w http.ResponseWriter) {
	writeJSON(w, domain.AccessKeyPermissions{
		Key:         "test-key",
		Permissions: &domain.Permissions{Send: []string{"event1"}, Receive: []string{"event2"}},
	})
}

func handleUpdateServiceKeyPair(w http.ResponseWriter) {
	writeJSON(w, domain.ServiceKeyPair{
		PublicKey:  "new-public-key",
		PrivateKey: "new-private-key",
	})
}

func handleListWorkspaces(w http.ResponseWriter) {
	writeJSON(w, domain.WorkspaceList{
		ResultsLength: 1,
		Results: []*domain.Workspace{
			{ID: "ws-1", Name: "gms", Path: "gms"},
		},
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
