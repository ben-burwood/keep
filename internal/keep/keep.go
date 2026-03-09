package keep

import (
	"github.com/google/uuid"
)

type KeepUUID string

type Keep struct {
	UUID    KeepUUID `json:"uuid"`
	Content string   `json:"content"`
}

// NewKeep - Constructor to create a new Keep with a unique UUID and current time as created_at.
func NewKeep(content string) *Keep {
	return &Keep{
		UUID:    KeepUUID(uuid.NewString()),
		Content: content,
	}
}
