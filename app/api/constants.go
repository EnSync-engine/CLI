package api

const (
	// HTTP headers
	headerAccessKey   = "X-ACCESS-KEY"
	headerContentType = "Content-Type"
	headerAccept      = "Accept"
	contentTypeJSON   = "application/json"

	// API paths
	pathEvent                = "/event"
	pathAccessKey            = "/access-key"
	pathWorkspace            = "/workspace"
	pathServiceKeyPair       = "/access/service-key-pair"
	pathAccessKeyPermissions = "/access-key/%s/permissions"
	pathAccessKeyByID        = "/access-key/%s"
	pathEventByName          = "/event/%s"
	pathEventByID            = "/event/%s"
)
