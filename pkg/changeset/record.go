package changeset

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
)

// New starts an empty change set for an org.
func New(orgID string) *ChangeSet {
	var b [4]byte
	_, _ = rand.Read(b[:])
	now := time.Now().UTC()
	return &ChangeSet{
		Format: Format,
		Metadata: Metadata{
			ID:        "cs_" + now.Format("20060102T150405Z") + "_" + hex.EncodeToString(b[:]),
			CreatedAt: now.Format(time.RFC3339),
		},
		Scope:      Scope{OrgID: orgID},
		Operations: []Operation{},
	}
}

// NewIDPrefix marks a placeholder ID for an object a change set creates. A file
// names such an object with a {"$ref": "<handle>"} value instead, so a target
// parameter never holds one.
const NewIDPrefix = "new:"

// HandleOf returns the handle in a "new:<handle>" placeholder ID.
func HandleOf(id string) (string, bool) {
	if !strings.HasPrefix(id, NewIDPrefix) {
		return "", false
	}
	return strings.TrimPrefix(id, NewIDPrefix), true
}
