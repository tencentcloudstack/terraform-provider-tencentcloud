package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const idSeparator = "#"

// IdFormat joins the given segments into a composite resource id with the
// repo-standard "#" separator. The segment order is ALWAYS parent-first
// (largest to smallest): in "a#b#c", a is the parent of b and b is the
// parent of c (e.g. "instance_id#disk_id").
func IdFormat(s ...string) string {
	return strings.Join(s, idSeparator)
}

// IdParse splits a composite resource id joined by "#" into its segments.
// Prefer IdParseErr, which also validates the segment count and emptiness.
func IdParse(s string) []string {
	return strings.Split(s, idSeparator)
}

// IdParseErr splits a composite resource id joined by "#" into exactly
// segments non-empty parts. It returns an error when the id does not have
// the expected shape: wrong segment count or any empty segment. Callers
// should wrap the error with resource-specific context, e.g. "Makers
// project id must be in the format `zone_id#project_id`".
func IdParseErr(id string, segments int) ([]string, error) {
	parts := strings.Split(id, idSeparator)
	if len(parts) != segments {
		return nil, fmt.Errorf("invalid id %q, expected %d segments joined by %q, got %d", id, segments, idSeparator, len(parts))
	}
	for _, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("invalid id %q, expected %d non-empty segments joined by %q", id, segments, idSeparator)
		}
	}
	return parts, nil
}

// BuildToken generates the random state id used by data sources: their id
// carries no cloud meaning, so a random token is minted at read time
// (mirrors internal/helper.BuildToken for the framework side).
func BuildToken() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	return base64.StdEncoding.EncodeToString(buf)
}

// BuildUUID generates a fresh UUID string.
func BuildUUID() string {
	return uuid.New().String()
}
