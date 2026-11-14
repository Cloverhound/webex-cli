package changeset

import _ "embed"

// Schema is the JSON Schema for a v1 change set.
//
//go:embed schema.json
var Schema []byte
