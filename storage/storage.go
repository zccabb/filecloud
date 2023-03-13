package storage

import (
	"filecloud/auth"
	"filecloud/settings"
	"filecloud/share"
	"filecloud/users"
)

// Storage is a storage powered by a Backend which makes the necessary
// verifications when fetching and saving data to ensure consistency.
type Storage struct {
	Users    users.Store
	Share    *share.Storage
	Auth     *auth.Storage
	Settings *settings.Storage
}
