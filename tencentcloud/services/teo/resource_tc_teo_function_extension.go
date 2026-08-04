package teo

import "strings"

// ParseTeoFunctionOriginalName extracts the original function name from the
// concatenated name returned by DescribeFunctions API.
// The API returns name in format: originalName + "-" + zoneId + "-" + appId,
// where zoneId has the format "zone-xxxxxxx". We locate the "-zone-" substring
// to determine the boundary of the original name.
// If "-zone-" is not found, the original value is returned as-is.
func ParseTeoFunctionOriginalName(name string) string {
	idx := strings.Index(name, "-zone-")
	if idx < 0 {
		return name
	}
	return name[:idx]
}
