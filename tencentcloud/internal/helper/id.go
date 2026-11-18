package helper

import (
	"fmt"
	"strings"
)

const connect = "#"

func IdFormat(s ...string) string {
	return strings.Join(s, connect)
}

func IdParse(s string) []string {
	return strings.Split(s, connect)
}

// IdParseErr splits a composite resource id joined by "#" into exactly
// segments non-empty parts. It returns an error when the id does not have
// the expected shape: wrong segment count or any empty segment. Callers
// should wrap the error with resource-specific context, e.g. "Makers
// project id must be in the format `zone_id#project_id`".
func IdParseErr(id string, segments int) ([]string, error) {
	parts := strings.Split(id, connect)
	if len(parts) != segments {
		return nil, fmt.Errorf("invalid id %q, expected %d segments joined by %q, got %d", id, segments, connect, len(parts))
	}
	for _, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("invalid id %q, expected %d non-empty segments joined by %q", id, segments, connect)
		}
	}
	return parts, nil
}
